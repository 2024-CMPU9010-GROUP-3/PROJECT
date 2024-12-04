import torch
from ultralytics import YOLO
import cv2
import os
import numpy as np
import math
import csv
from geopy.distance import geodesic
from shapely.geometry import box


def create_mask(image_path, save_path, threshold=240):
    """
    Creates and saves a binary mask from the mapbox image of the road (Mapbox Streets).
    The roads are in white and some additional roads are in orange/yellow (highways/ roads with more lanes).
    Most of the street names are removed in a way as to not distort the road sizes. (Some text remains as then 
    smaller roads (where the name takes up the entire width) would be completly erased)

    Params:
        image_path (str): Path of the image
        save_path (str): Path to save the mask
        threshold (int): Threshold to differentiate the road from the areas outside of the road
    """
    img = cv2.imread(image_path)
    img_gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    _, road_mask = cv2.threshold(img_gray, threshold, 255, cv2.THRESH_BINARY)

    img_hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    lower_orange = np.array([10, 100, 100])
    upper_orange = np.array([25, 255, 255])
    lower_yellow = np.array([25, 100, 100])
    upper_yellow = np.array([35, 255, 255])

    orange_mask = cv2.inRange(img_hsv, lower_orange, upper_orange)
    yellow_mask = cv2.inRange(img_hsv, lower_yellow, upper_yellow)

    dilation_kernel = np.ones((15, 15), np.uint8)#we thicken the road width for highways as the road doesn't take into account the multiple lanes (to reduce misclassifications)
    orange_mask_dilated = cv2.dilate(orange_mask, dilation_kernel, iterations=2)
    yellow_mask_dilated = cv2.dilate(yellow_mask, dilation_kernel, iterations=2)

    combined_mask = cv2.bitwise_or(road_mask, orange_mask_dilated)
    combined_mask = cv2.bitwise_or(combined_mask, yellow_mask_dilated)

    kernel = np.ones((2, 2), np.uint8)#use smaller kernel as it works better
    combined_mask = cv2.morphologyEx(combined_mask, cv2.MORPH_CLOSE, kernel)#cv2.MORPH_CLOSE actually works better

    contours, _ = cv2.findContours(combined_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    mask_filtered = np.zeros_like(combined_mask)

    for contour in contours:
        if cv2.contourArea(contour) > 300: 
            cv2.drawContours(mask_filtered, [contour], -1, 255, thickness=cv2.FILLED)

    cv2.imwrite(save_path, mask_filtered)


def detect_parking_spots_in_image(image_path, road_mask_path, output_image_path, model, conf_threshold=0.4):
    """
    Detect cars in the image using the retrained YOLO model and draws them on the image.

    Params:
        image_path (str): Path of the image
        road_mask_path (str): Path of the saved mask
        output_image_path (str): Path to save the image with bounding boxes, red for parking and blue cars on the road
        model : YOLO model
        conf_threshold (float): Minimum confidence threshold for a predictions to be considered
        
    Returns:
        detections_parking (list): List of the bounding boxes for the cars not on the road [x_center, y_center, width, height, angle_degrees]
    """
    img = cv2.imread(image_path)
    results = model([img])

    detections_parking = []

    road_mask = cv2.imread(road_mask_path, cv2.IMREAD_GRAYSCALE)

    if not results:
        print("No cars detected.")
        return detections_parking

    for result in results:
        detections = result.obb

        for box in detections:
            x_center, y_center, width, height, angle_radians = map(float, box.xywhr[0])
            angle_degrees = angle_radians * (180 / math.pi) 
            cls = int(box.cls[0])
            conf = float(box.conf[0])

            if conf < conf_threshold:#we ignore the predictions that have a low confidance score as they are more likely to be misclassifications
                continue

            if cls == 0:
                if -45 <= angle_degrees <= 45 or 135 <= angle_degrees <= 225:
                    orientation = "horizontal"
                else:
                    orientation = "vertical"
                    
                rect = ((x_center, y_center), (width, height), angle_degrees)
                box_points = cv2.boxPoints(rect)
                box_points = np.int32(box_points)

                x_min = int(x_center - width / 2)
                x_max = int(x_center + width / 2)
                y_min = int(y_center - height / 2)
                y_max = int(y_center + height / 2)

                car_region = road_mask[y_min:y_max, x_min:x_max]

                if car_region.size == 0:
                    continue

                road_pixels = cv2.countNonZero(car_region)
                total_pixels = car_region.size

                if road_pixels / total_pixels > 0.5:
                    print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is on the road")
                    detections_parking.append([x_center, y_center, width, height, 1])
                    cv2.polylines(img, [box_points], isClosed=True, color=(255, 0, 0), thickness=2) #blue if on the road
                else:
                    print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is not on the road (possibly parked)")
                    detections_parking.append([x_center, y_center, width, height, 0])
                    cv2.polylines(img, [box_points], isClosed=True, color=(0, 0, 255), thickness=2) #red if parked

    cv2.imwrite(output_image_path, img)
    return detections_parking

def get_predictions_in_image(model, longitude, latitude, directory):
    """
    Returns all empty parking spot predictions in an image in the correct format for evaluation

    Params:
        model : YOLO model
        longitude (float): Longitude value of the image
        latitude (float): Latitude value of the image
        directory(str): Path to the directory containing the images and the labels in json format 


    Returns:
            empty_detections (list): List of all the empty parking spots found in the image in the format x, y, width, height, orientation (x and y being pixel values)
    """
    output_folder = directory
    output_path_satelite_image = os.path.join(output_folder, f'{longitude}_{latitude}_satellite.png')
    output_path_road_image = os.path.join(output_folder, f'{longitude}_{latitude}_road.png')
    output_path_mask_image = os.path.join(output_folder, f'{longitude}_{latitude}_mask.png')
    output_path_bb_image = os.path.join(output_folder, f'{longitude}_{latitude}_bounding_boxes.png')

    create_mask(output_path_road_image, output_path_mask_image)
    detections = detect_parking_spots_in_image(output_path_satelite_image, output_path_mask_image, output_path_bb_image, model)

    return detections


def get_true_labels(long, lat, directory, image_width=400, image_height=400):
    """
    Retrieve the true labels for a specific image

    Params:
        long (float): Longitude of the image
        lat (float): Latitude of the image
        directory(str): Path to the directory containing the images and the labels in a txt file in the YOLO format 
        image_width (int): Image width in pixels
        image_height (int): Image height in pixels

    Returns:
        true_labels (list): List of true labels bounding boxes in the format x_pixel, y_pixel, width, height, class
    """
    true_labels_file = f"{long}_{lat}_satellite.txt"
    file_path = os.path.join(directory, true_labels_file)

    if not os.path.exists(file_path):
        print(f"Error: Label file {true_labels_file} does not exist in {directory}.")
        return []

    true_labels = []

    try:
        with open(file_path, 'r') as file:
            for line in file:
                parts = line.strip().split()
                
                if len(parts) != 5:
                    print(f"Warning: Skipping line due to unexpected format: {line}")
                    continue
                
                try:
                    classification = float(parts[0])
                    x_pixel = float(parts[1])*image_width #the values given by label studio are normalized and we want the denormalized values to compare with the predictions
                    y_pixel = float(parts[2])*image_height
                    width = float(parts[3])*image_width
                    height = float(parts[4])*image_height
                    true_labels.append([x_pixel, y_pixel, width, height, classification])
                except ValueError:
                    print(f"Warning: Invalid data format in line: {line}")
                    continue
    except IOError as e:
        print(f"Error reading file {true_labels_file}: {e}")

    return true_labels

def evaluate_predictions(predictions, true_labels, iou_threshold=0.35):
    """
    Evaluates our predictions in an image with the true labels
    Returns Average IoU, Precision, Recall, F1 Score, Accuracy, Specificity per class and the overall Balanced Accuracy.

    Params:
        predictions (list): List of predictions
        true_labels (list): List of true_labels
        iou_threshold (float): IoU threshold to consider a prediction correct

    Returns:
        avg_iou, precision_parked, recall_parked, f1_score_parked, accuracy_parked, specificity_parked, precision_road, recall_road, f1_score_road, accuracy_road, specificity_road, balanced_accuracy (float): Metrics per class and overall balanced average and average iou
    """

    def calculate_iou(box1, box2):
        """Calculates IoU for two bounding boxes.
        Box1 is the predictions and box2 is the true label"""
        x1, y1, w1, h1 = box1
        x2, y2, w2, h2 = box2

        box1_tl = (x1 - w1 / 2, y1 - h1 / 2)
        box1_br = (x1 + w1 / 2, y1 + h1 / 2)
        box2_tl = (x2 - w2 / 2, y2 - h2 / 2)
        box2_br = (x2 + w2 / 2, y2 + h2 / 2)

        x_left = max(box1_tl[0], box2_tl[0])
        y_top = max(box1_tl[1], box2_tl[1])
        x_right = min(box1_br[0], box2_br[0])
        y_bottom = min(box1_br[1], box2_br[1])

        if x_right < x_left or y_bottom < y_top:
            return 0  #no overlap

        intersection_area = (x_right - x_left) * (y_bottom - y_top)

        box1_area = w1 * h1
        box2_area = w2 * h2
        union_area = box1_area + box2_area - intersection_area

        return intersection_area / union_area

    true_positives = {0: 0, 1: 0}
    false_positives = {0: 0, 1: 0}
    false_negatives = {0: 0, 1: 0}
    true_negatives = {0: 0, 1: 0}
    iou_scores = []

    matched_labels = set()

    # For each prediction, find the closest label and compare
    for pred in predictions:
        best_iou = 0
        best_match = None
        pred_box, pred_class = pred[:4], pred[4]
        for i, gt in enumerate(true_labels):
            gt_box, gt_class = gt[:4], gt[4]
            iou = calculate_iou(pred_box, gt_box)
            if iou > best_iou and pred_class == gt_class:
                best_iou = iou
                best_match = i

        if best_iou >= iou_threshold:
            if best_match not in matched_labels:
                true_positives[pred_class] += 1
                iou_scores.append(best_iou)
                matched_labels.add(best_match)
        else:
            false_positives[pred_class] += 1

    for i, gt in enumerate(true_labels):
        gt_class = gt[4]
        if i not in matched_labels:
            false_negatives[gt_class] += 1

    for cls in [0, 1]:
        true_negatives[cls] = sum(len(true_labels) - true_positives[c] - false_negatives[c] - false_positives[c] for c in [0, 1] if c != cls)

    metrics = {}
    for cls in [0, 1]: 
        tp = true_positives[cls]
        fp = false_positives[cls]
        fn = false_negatives[cls]
        tn = true_negatives[cls]

        precision = tp / (tp + fp) if (tp + fp) > 0 else (None if tp == 0 and fp == 0 else 0)
        recall = tp / (tp + fn) if (tp + fn) > 0 else (None if tp == 0 and fn == 0 else 0)
        f1_score = 2 * (precision * recall) / (precision + recall) if precision is not None and recall is not None and (precision + recall) > 0 else 0
        accuracy = (tp + tn) / (tp + fp + fn + tn) if (tp + fp + fn + tn) > 0 else (None if tp + tn == 0 and fp + fn == 0 else 0)
        specificity = tn / (tn + fp) if (tn + fp) > 0 else (None if tn == 0 and fp == 0 else 0)

        metrics[cls] = (precision, recall, f1_score, accuracy, specificity)

    balanced_accuracy = sum((metrics[cls][1] if metrics[cls][1] is not None else 0) + (metrics[cls][4] if metrics[cls][4] is not None else 0) for cls in [0, 1]) / 4
    avg_iou = sum(iou_scores) / len(iou_scores) if iou_scores else 0

    return (avg_iou,
        metrics[0][0], metrics[0][1], metrics[0][2], metrics[0][3], metrics[0][4],  #Parked car metrics
        metrics[1][0], metrics[1][1], metrics[1][2], metrics[1][3], metrics[1][4],  #Car on road metrics
        balanced_accuracy)

def main(directory, output_file="metrics_road_mask.csv"):
    """
    Main function to evaluate the classification of cars into on the road vs parked on the test images and calculate the corresponding performance metrics
    Saves all metrics in a csv file for each image as well as the overall average metrics.

    Params:
        directory(str): Path to the directory containing the images and the labels in json format 
        output_file (str): Path to save the CSV file ctaining the metrics for each image and the overall metrics

    """
    if not os.path.exists(directory):
        print(f"Error: Directory {directory} does not exist.")
        return
    
    files = os.listdir(directory)

    coordinates = []
    for file in files:
        if file.endswith(".txt"):
            try:
                long, lat, _ = file.split("_")
                long, lat = float(long), float(lat)
                coordinates.append((long, lat))
            except ValueError:
                print(f"Skipping file {file}: Invalid name format for extracting coordinates.")
                continue
    
    model = YOLO("best - obb.pt")

    metrics = {
        "iou": [],
        "precision_parked": [],
        "precision_road": [],
        "recall_parked": [],
        "recall_road": [],
        "f1_score_parked": [],
        "f1_score_road": [],
        "accuracy_parked": [],
        "accuracy_road": [],
        "specificity_parked": [],
        "specificity_road": [],
        "balanced_accuracy": []}

    image_metrics = []

    for long, lat in set(coordinates):
        predictions = get_predictions_in_image(model, long, lat, directory)
        true_labels = get_true_labels(long, lat, directory)
        if not predictions and not true_labels:# if there are no detections and no true labels we want to skip the calculation of the metrics
            continue
        #print(predictions)
        #print(true_labels)

        iou, precision_parked, recall_parked, f1_score_parked, accuracy_parked, specificity_parked, precision_road, recall_road, f1_score_road, accuracy_road, specificity_road, balanced_accuracy  = evaluate_predictions(predictions, true_labels)

        metrics["iou"].append(iou)
        metrics["precision_parked"].append(precision_parked)
        metrics["precision_road"].append(precision_road)
        metrics["recall_parked"].append(recall_parked)
        metrics["recall_road"].append(recall_road)
        metrics["f1_score_parked"].append(f1_score_parked)
        metrics["f1_score_road"].append(f1_score_road)
        metrics["accuracy_parked"].append(accuracy_parked)
        metrics["accuracy_road"].append(accuracy_road)
        metrics["specificity_parked"].append(specificity_parked)
        metrics["specificity_road"].append(specificity_road)
        metrics["balanced_accuracy"].append(balanced_accuracy)

        image_metrics.append({
            "longitude": long,
            "latitude": lat,
            "iou": iou,
            "precision_parked": precision_parked,
            "precision_road": precision_road,
            "recall_parked": recall_parked,
            "recall_road": recall_road,
            "f1_score_parked": f1_score_parked,
            "f1_score_road": f1_score_road,
            "accuracy_parked": accuracy_parked,
            "accuracy_road": accuracy_road,
            "specificity_parked": specificity_parked,
            "specificity_road": specificity_road,
            "balanced_accuracy": balanced_accuracy,
        })

        print(f"Metrics for image {long}, {lat}: IoU={iou}, Precision Parked={precision_parked}, Precision Road={precision_road}, Recall Parked={recall_parked}, Recall Road={recall_road}, F1 Score Parked={f1_score_parked}, F1 Score Road={f1_score_road}, Accuracy Parked={accuracy_parked}, Accuracy Road={accuracy_road}, Specificity Parked={specificity_parked}, Specificity Road={specificity_road},  Balanced Accuracy ={balanced_accuracy}")

    overall_metrics = {key: np.nanmean([value for value in values if value is not None]) if len(values) > 0 else None for key, values in metrics.items()}
    overall_metrics["longitude"] = "Overall"
    overall_metrics["latitude"] = "Metrics"

    print("Overall Metrics:")
    for key, value in overall_metrics.items():
        if key not in {"longitude", "latitude"}:
            print(f"Average {key.replace('_', ' ').title()}: {value:.2f}")

    with open(output_file, mode="w", newline="") as file:
        writer = csv.DictWriter(file, fieldnames=[
            "longitude", "latitude", "iou", "precision_parked", "precision_road", "recall_parked", "recall_road",
            "f1_score_parked", "f1_score_road", "accuracy_parked", "accuracy_road", "specificity_parked",
            "specificity_road", "balanced_accuracy"
        ])
        writer.writeheader()
        writer.writerows(image_metrics)
        writer.writerow(overall_metrics)

    print(f"Metrics saved to {output_file}")

if __name__ == "__main__":
    main("all_test_images_road")