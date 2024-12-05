import torch
from ultralytics import YOLO
import cv2
import os
import numpy as np
import math
import pandas as pd
import random
import csv
from geopy.distance import geodesic
from sklearn.cluster import DBSCAN

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
                else:
                    print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is not on the road (possibly parked)")
                    detections_parking.append([x_center, y_center, width, height, angle_degrees, orientation])

    cv2.imwrite(output_image_path, img)
    return detections_parking

def convert_bounding_box_to_coordinates(x, y, longitude, latitude):
    """
    Convert pixel coordinates of the bounding box center to longitude and latitude
    
    Params:
        x, y (float): Pixel coordinates
        longitude (float): Longitude of the image center
        latitude (float): Latitude of the image center
    
    Returns:
        long, lat (float): The longitude and latitude of the center of the bounding box
    """
    num_tiles = 2 ** 18
    tile_size = 256

    lat_rad = math.radians(latitude)

    center_x_tile = (longitude + 180.0) / 360.0 * num_tiles
    center_y_tile = (1.0 - math.log(math.tan(lat_rad) + (1 / math.cos(lat_rad))) / math.pi) / 2.0 * num_tiles

    center_tile_x = center_x_tile * tile_size
    center_tile_y = center_y_tile * tile_size

    meters_per_pixel = 156543.03392 * math.cos(lat_rad) / (2 ** 18)

    pixel_x_offset = (x - 200) * meters_per_pixel
    pixel_y_offset = (y - 200) * meters_per_pixel

    long = (center_tile_x + pixel_x_offset) / (num_tiles * tile_size) * 360.0 - 180.0
    
    lat_rad = math.atan(math.sinh(math.pi * (1 - 2 * (center_tile_y + pixel_y_offset) / (num_tiles * tile_size))))
    lat = math.degrees(lat_rad)

    return long, lat

def convert_coordinates_to_bounding_box(longitude, latitude, center_long, center_lat):
    """
    Convert longitude and latitude to pixel coordinates
    
    Params:
        longitude (float): Longitude of the object (car or empty parking space)
        latitude (float): Latitude of the object (car or empty parking space)
        center_long (float): Longitude of the center of the image
        center_lat (float): Latitude of the center of the image.
        
    Returns:
        x, y (float): pixel coordinates of the center of the bounding box
    """
    num_tiles = 2 ** 18
    tile_size = 256

    center_lat_rad = math.radians(center_lat)
    center_x_tile = (center_long + 180.0) / 360.0 * num_tiles
    center_y_tile = (1.0 - math.log(math.tan(center_lat_rad) + (1 / math.cos(center_lat_rad))) / math.pi) / 2.0 * num_tiles

    lat_rad = math.radians(latitude)
    x_tile = (longitude + 180.0) / 360.0 * num_tiles
    y_tile = (1.0 - math.log(math.tan(lat_rad) + (1 / math.cos(lat_rad))) / math.pi) / 2.0 * num_tiles

    center_x_pixel = center_x_tile * tile_size
    center_y_pixel = center_y_tile * tile_size
    x_pixel = x_tile * tile_size
    y_pixel = y_tile * tile_size

    meters_per_pixel = 156543.03392 * math.cos(center_lat_rad) / num_tiles
    x_offset = (x_pixel - center_x_pixel) / meters_per_pixel + 400 / 2
    y_offset = (y_pixel - center_y_pixel) / meters_per_pixel + 400 / 2

    return x_offset, y_offset


def calculate_avg_spot_dimensions(cars):
    """
    Calculates average parking spot width and length in pixels
    
    Params:
        cars (list): List of corrdinates with width and length
        
    Returns:
        avg_width_meters, avg_length_meters, avg_width_pixels, avg_length_pixels (float): Average width and length of the cars identified in meters and in pixels
    """
    widths = [car[2] for car in cars]
    lengths = [car[3] for car in cars]
    avg_width_pixels = np.median(widths)
    avg_length_pixels = np.median(lengths)

    #the avgerage width and length are set as in practise it works better and accounts for the variations and avoids misclassifications when avg_length_meters or avg_width_meters are < 3 (when they were calculated dynamically)
    avg_width_meters = 3.05
    avg_length_meters = 3.05

    return avg_width_meters, avg_length_meters, avg_width_pixels, avg_length_pixels

def detect_empty_spots(cars, avg_spot_width, avg_spot_length, gap_threshold_meters=12, duplicate_threshold_meters=1, overlap_threshold_meters=1.25):
    """
    Detects empty spots in a group of parked cars in an image, based on the detected car bounding box centers.
    There are 4 cases: horizontal in a row, horizontal stacked, vertical in a column and vertical side by side.
    Duplicate spots and spots coinciding with cars identified by the Yolo model are removed.
    
    Params:
        cars (list): List of car bounding box centers with orientation horizontal or vertical
        avg_spot_width (float): Average width of a parking spot in meters
        avg_spot_length (float): Average length of a parking spot in meters
        gap_threshold_meters (float): Maximum allowed gap to consider there is an empty parking spot or multiple parking spots
        duplicate_threshold_meters (float): Threshold to differenciate between spots that are considered identical in meters
        overlap_threshold_meters (float): Threshold to remove empty spots overlapping with detected cars

    Returns:
       empty_spots (list): List of coordinates of estimated empty parking spots with horizontal or vertical orientation
    """
    empty_spots = []
    
    horizontal_cars_sorted_by_lat = sorted([car for car in cars if car[5] == 'horizontal'], key=lambda point: point[1])
    vertical_cars_sorted_by_long = sorted([car for car in cars if car[5] == 'vertical'], key=lambda point: point[0])  

    def find_empty_spots(sorted_cars, alignment, gap_dimension, gap_threshold_meters):
        """ Detects empty spots in the sorted list of cars for a specific alignment (horizontal or vertical) 
        
        Params:
        sorted_cars (list): List of car bounding box centers sorted by orientation (horizontal/vertical) and coordinates (long/lat)
        alignment (str): Either horizontal or vertical
        gap_dimension (str): Either avg_spot_width or avg_spot_length
        gap_threshold_meters (float): Maximum allowed gap to consider there is an empty parking spot or multiple parking spots
        """
        for i in range(len(sorted_cars) - 1):
            x_current, y_current, _, _, angle_current, _ = sorted_cars[i]
            x_next, y_next, _, _, angle_next, _ = sorted_cars[i + 1]

            gap_distance = geodesic((y_current, x_current), (y_next, x_next)).meters
            avg_half_dim = gap_dimension / 2
            adjusted_gap = gap_distance - 2 * avg_half_dim

            angle_radians = math.atan2(y_next - y_current, x_next - x_current)
            angle_degrees = math.degrees(angle_radians)

            angle_deviation = abs(angle_current - angle_next)
            angle_deviation = min(angle_deviation, 360 - angle_deviation)

            if angle_deviation > 35:
                continue
            
            if adjusted_gap <= gap_threshold_meters and adjusted_gap > gap_dimension:
                num_spots = int(adjusted_gap // gap_dimension)
                
                for j in range(1, num_spots + 1):
                    empty_x_center = x_current + j * (x_next - x_current) / (num_spots + 1)
                    empty_y_center = y_current + j * (y_next - y_current) / (num_spots + 1)
                    empty_spots.append(([empty_x_center, empty_y_center], angle_degrees, alignment))
                    print(f"Empty parking spot at {empty_x_center}, {empty_y_center}")

    find_empty_spots(horizontal_cars_sorted_by_lat, 'horizontal', avg_spot_width, gap_threshold_meters) #Horizontal spots in a row
    find_empty_spots(vertical_cars_sorted_by_long, 'vertical', avg_spot_width, gap_threshold_meters)  #Vertical spots in columns

    empty_spots = sorted(empty_spots, key=lambda spot: (spot[0][0], spot[0][1]))
    unique_empty_spots = []

    for i, spot in enumerate(empty_spots):
        if i == 0:
            unique_empty_spots.append(spot)
        else:
            distance_to_prev = geodesic((spot[0][1], spot[0][0]), (empty_spots[i - 1][0][1], empty_spots[i - 1][0][0])).meters
            
            if distance_to_prev >= duplicate_threshold_meters or spot[2] != empty_spots[i - 1][2]:
                unique_empty_spots.append(spot)
            else:
                print(f"Removed spot at {spot[0]} due to proximity to {empty_spots[i - 1][0]}. Distance: {distance_to_prev:.2f}")


    filtered_empty_spots = []

    for empty_spot in unique_empty_spots:
        empty_x, empty_y = empty_spot[0]
        overlap = False
        for car in cars:
            car_x, car_y, _, _, _, _ = car
            distance_to_car = geodesic((empty_y, empty_x), (car_y, car_x)).meters
            if distance_to_car < overlap_threshold_meters:
                overlap = True
                break
        if not overlap:
            filtered_empty_spots.append(empty_spot)

    #print(f'Empty spots: {len(empty_spots)}')
    #print(f'Unique empty spots: {len(unique_empty_spots)}')
    #print(f'Filtered empty spots: {len(filtered_empty_spots)}')

    return filtered_empty_spots

def filter_empty_spots_on_road(empty_spots, road_mask_path, center_long, center_lat, avg_spot_width, avg_spot_length):
    """
    Filters out empty parking spots found that are on the road

    Params:
        empty_spots (list): List of coordinates of estimated empty parking spots with horizontal or vertical orientation
        road_mask_path (ndarray): Path of the saved mask
        center_long (float): Longitude of the center of the image
        center_lat (float): Latitude of the center of the image
        avg_spot_width (float): Average width of a parking spot in pixels
        avg_spot_length (float): Average length of a parking spot in pixels

    Returns:
        filtered_empty_spots (list): List of empty parking spots not on the road
    """
    road_mask = cv2.imread(road_mask_path, cv2.IMREAD_GRAYSCALE)

    filtered_empty_spots = []
    
    for spot, rotation, alignment in empty_spots:
        x_center, y_center = convert_coordinates_to_bounding_box(spot[0], spot[1], center_long, center_lat)

        if alignment == 'horizontal':
            x_min = int(x_center - avg_spot_width / 2)
            x_max = int(x_center + avg_spot_width / 2)
            y_min = int(y_center - avg_spot_length / 2)
            y_max = int(y_center + avg_spot_length / 2)
        else:
            x_min = int(x_center - avg_spot_length / 2)
            x_max = int(x_center + avg_spot_length / 2)
            y_min = int(y_center - avg_spot_width / 2)
            y_max = int(y_center + avg_spot_width / 2)

        car_region = road_mask[y_min:y_max, x_min:x_max]
        road_pixels = cv2.countNonZero(car_region)
        total_pixels = car_region.size

        if total_pixels > 0 and road_pixels / total_pixels <= 0.4:
            filtered_empty_spots.append((spot[0], spot[1], avg_spot_width, avg_spot_length, rotation, alignment))

    return filtered_empty_spots


def classify_parking_spots(all_parking_spots, road_mask_path, center_long, center_lat, road_proximity_threshold=30, parking_lot_min_spots=18, clustering_eps=55, clustering_min_samples=5):
    """
    Classifies parking spots as public(on the street parking), private(residential) or parking lot based on their proximity to the road (calculated using the road mask).
    Parking lots are identified through clustering using DBSCAN

    Params:
        all_parking_spots (list): List of all parking spots(by the model and then the empty parking detection) 
        road_mask_path (string): Path to road mask
        center_long (float): Longitude of the center of the image
        center_lat (float): Latitude of the center of the image
        road_proximity_threshold (int): Threshold in pixels to classify a spot near the road as public
        parking_lot_min_spots (int): Minimum number of spots in a cluster to classify it as a parking lot
        clustering_eps (float): Maximum distance between spots in pixels to form a cluster
        clustering_min_samples (int): Minimum number of samples to form a cluster

    Returns:
        classified_spots (list): List of parking spots with classification added
    """
    classified_spots = []

    road_mask = cv2.imread(road_mask_path, cv2.IMREAD_GRAYSCALE)
    road_contours, _ = cv2.findContours(road_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    #print(center_long, center_lat)
    pixel_coords = []

    for long, lat, _, _, _, _ in all_parking_spots:
        x_center, y_center = convert_coordinates_to_bounding_box(long, lat, center_long, center_lat)
        pixel_coords.append([x_center, y_center])

    clustering = DBSCAN(eps=clustering_eps, min_samples=clustering_min_samples).fit(pixel_coords)
    labels = clustering.labels_

    for idx, spot in enumerate(all_parking_spots):
        x_center, y_center = pixel_coords[idx]
        cluster_label = labels[idx]

        min_distance = float("inf")

        for contour in road_contours:
            distance = cv2.pointPolygonTest(contour, (x_center, y_center), measureDist=True)
            min_distance = min(min_distance, abs(distance))

        classification = 1 #private

        if min_distance <= road_proximity_threshold:
            classification = 2 #public

        if cluster_label != -1:
            cluster_size = np.sum(labels == cluster_label)
            if cluster_size >= parking_lot_min_spots:
                classification = 0 #parking lot

        classified_spots.append([pixel_coords[idx][0], pixel_coords[idx][1], spot[2], spot[3], spot[5], classification])
        #print(classification)

    return classified_spots

def draw_classification(image_path, spots):
    """
    Draws cluster labels and classifications labels on the image for each spot.

    Params:
        image_path (str): Path of the image
        spots (list): List of parking spot coordinates and classification labels (public, private or parking lot)
        cluster_labels (list): List of cluster labels corresponding to each spot
        center_long (float): Longitude of the center of the image
        center_lat (float): Latitude of the center of the image
    """
    image = cv2.imread(image_path)

    for x_pixel, y_pixel, width, height, orientation, classification in spots:

        if orientation == "horizontal":
            x1 = int(x_pixel - width // 2)
            y1 = int(y_pixel - height // 2)
            x2 = int(x_pixel + width // 2)            
            y2 = int(y_pixel + height // 2)
        else:
            x1 = int(x_pixel - height // 2)
            y1 = int(y_pixel - width // 2)
            x2 = int(x_pixel + height // 2)            
            y2 = int(y_pixel + width // 2)
            

        if classification == 1: #Draw private in red
            color = (0, 0, 255)
        elif classification == 2:  #Draw public in green
            color = (0, 255, 0)
        elif classification == 0: #Draw parking lot in blue
            color = (255, 0, 0) 

        cv2.rectangle(image, (x1, y1), (x2, y2), color, 2)

    cv2.imwrite(image_path, image)

def get_predictions_in_image(model, longitude, latitude, directory):
    """
    Detects all the parking spaces in the image (at longitude/latitude) and returns a list of coordinates, agregating
    the cars (of the road) found by the Yolo model and the empty parking spots found (which are drawn and added to the image) 

    Params:
        model : YOLO model
        longitude (float): Longitude value
        latitude (float): Latitude value
        directory(str): Path to the directory containing the images and the labels in a txt file in the YOLO format 

    Returns: 
        all_detections (list): List of all coordinates of parking spots found in the image in the format long, lat, classification
    """
    output_path_satelite_image = os.path.join(directory, f'{longitude}_{latitude}_satellite.png')
    output_path_road_image = os.path.join(directory, f'{longitude}_{latitude}_road.png')
    output_path_mask_image = os.path.join(directory, f'{longitude}_{latitude}_mask.png')
    output_path_bb_image = os.path.join(directory, f'{longitude}_{latitude}_bounding_boxes.png')

    create_mask(output_path_road_image, output_path_mask_image)
    detections = detect_parking_spots_in_image(output_path_satelite_image, output_path_mask_image, output_path_bb_image, model)

    all_detections = []

    for detection in detections:
        x_center, y_center, width, height, angle, orientation = detection
        long, lat = convert_bounding_box_to_coordinates(x_center, y_center, longitude, latitude)

        if isinstance(long, torch.Tensor):
            long = long.item()
        if isinstance(lat, torch.Tensor):
            lat = lat.item()

        print(f"Car coordinates: ({long}, {lat})")
        all_detections.append([long, lat, width, height, angle, orientation]) 

    if all_detections:
        avg_width_meters, avg_length_meters, avg_width_pixels, avg_length_pixels = calculate_avg_spot_dimensions(all_detections)
        empty_spots = detect_empty_spots(all_detections, avg_width_meters, avg_length_meters)
        empty_spots_filtered = filter_empty_spots_on_road(empty_spots, output_path_mask_image, longitude, latitude, avg_width_pixels, avg_length_pixels)
        #print(f'Filtered spots after road mask: {len(empty_spots_filtered)}')
        all_detections.extend(empty_spots_filtered)
        all_detections = classify_parking_spots(all_detections, output_path_mask_image, longitude, latitude)
        draw_classification(output_path_bb_image, all_detections)
        all_detections = [[x_pixel, y_pixel, width, height, classification] for x_pixel, y_pixel, width, height, _, classification in all_detections]
    return all_detections

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
        true_labels (list): List of true labels bounding boxes in the format x_pixel, y_pixel, width, height, orientation
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
                    x_pixel = float(parts[1])*image_width #the values given by label studio are normalized and we want the denormalized values to compare with the predictions
                    y_pixel = float(parts[2])*image_height
                    width = float(parts[3])*image_width
                    height = float(parts[4])*image_height
                    classification = int(parts[0])
                    true_labels.append([x_pixel, y_pixel, width, height, classification])
                except ValueError:
                    print(f"Warning: Invalid data format in line: {line}")
                    continue
    except IOError as e:
        print(f"Error reading file {true_labels_file}: {e}")

    return true_labels

def draw_true_labels(true_labels, directory, longitude, latitude):
    """
    Draws true labels to visualize the evaluation

    Params:
        true_labels (list): List of true labels bounding boxes in the format x_pixel, y_pixel, width, height, orientation
        directory(str): Path to the directory containing the images and the labels in a txt file in the YOLO format 
        longitude (float): Longitude of the image
        latitude (float): Latitude of the image
    """
    image_path = os.path.join(directory, f'{longitude}_{latitude}_bounding_boxes.png')
    image = cv2.imread(image_path)

    for x_pixel, y_pixel, width, height, classification in true_labels:
        
        x1 = int(x_pixel - width // 2)
        y1 = int(y_pixel - height // 2)
        x2 = int(x_pixel + width // 2)            
        y2 = int(y_pixel + height // 2)

        if classification == 1: #Draw private true labels in pink
            color = (180, 105, 255)
        elif classification == 2:  #Draw public true labels in teal
            color = (128, 128, 0)
        elif classification == 0: #Draw parking lot true label in cyan
            color = (255, 255, 0) 

        cv2.rectangle(image, (x1, y1), (x2, y2), color, 2)

    cv2.imwrite(image_path, image)

def evaluate_predictions(predictions, true_labels, iou_threshold=0.4):
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

    all_classes = set([p[4] for p in predictions] + [t[4] for t in true_labels])# we dynamically set the number of classes to handle the cases when not all classes are present in the image

    true_positives = {cls: 0 for cls in all_classes}
    false_positives = {cls: 0 for cls in all_classes}
    false_negatives = {cls: 0 for cls in all_classes}
    true_negatives = {cls: 0 for cls in all_classes}
    iou_scores = []

    matched_labels = set()

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

    for cls in all_classes:
        true_negatives[cls] = sum(
            len(true_labels) - true_positives[c] - false_negatives[c] - false_positives[c]
            for c in all_classes if c != cls
        )

    metrics_per_class = {}

    for cls in all_classes:
        tp = true_positives[cls]
        fp = false_positives[cls]
        fn = false_negatives[cls]
        tn = true_negatives[cls]

        precision = tp / (tp + fp) if (tp + fp) > 0 else (None if tp == 0 and fp == 0 else 0)
        recall = tp / (tp + fn) if (tp + fn) > 0 else (None if tp == 0 and fn == 0 else 0)
        f1_score = 2 * (precision * recall) / (precision + recall) if precision is not None and recall is not None and (precision + recall) > 0 else 0
        accuracy = (tp + tn) / (tp + fp + fn + tn) if (tp + fp + fn + tn) > 0 else (None if tp + tn == 0 and fp + fn == 0 else 0)
        specificity = tn / (tn + fp) if (tn + fp) > 0 else (None if tn == 0 and fp == 0 else 0)

        metrics_per_class[cls] = {
            "precision": precision,
            "recall": recall,
            "f1_score": f1_score,
            "accuracy": accuracy,
            "specificity": specificity
        }

    balanced_accuracy = sum(
        (metrics_per_class[cls]["recall"] if metrics_per_class[cls]["recall"] is not None else 0) +
        (metrics_per_class[cls]["specificity"] if metrics_per_class[cls]["specificity"] is not None else 0)
        for cls in all_classes
    ) / (2 * len(all_classes))

    avg_iou = sum(iou_scores) / len(iou_scores) if iou_scores else 0

    #print(metrics_per_class)

    return avg_iou, metrics_per_class, balanced_accuracy

def main(directory, output_file="metrics_spots_classification.csv"):
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
        "balanced_accuracy": [],
    }
    classes = [0, 1, 2]
    for cls in classes:
        metrics.update({
            f"precision_{cls}": [],
            f"recall_{cls}": [],
            f"f1_score_{cls}": [],
            f"accuracy_{cls}": [],
            f"specificity_{cls}": [],
        })

    image_metrics = []

    for long, lat in set(coordinates):
        predictions = get_predictions_in_image(model, long, lat, directory)
        true_labels = get_true_labels(long, lat, directory)
        draw_true_labels(true_labels, directory, long, lat)
        if not predictions and not true_labels:  #skip images if there are no detections and no true labels
            continue
        #print(true_labels)
        #print(predictions)

        avg_iou, per_class_metrics, balanced_accuracy = evaluate_predictions(predictions, true_labels)

        metrics["iou"].append(avg_iou)
        for cls in classes:
            if cls in per_class_metrics:
                cls_metrics = per_class_metrics[cls]
                metrics[f"precision_{cls}"].append(cls_metrics["precision"])
                metrics[f"recall_{cls}"].append(cls_metrics["recall"])
                metrics[f"f1_score_{cls}"].append(cls_metrics["f1_score"])
                metrics[f"accuracy_{cls}"].append(cls_metrics["accuracy"])
                metrics[f"specificity_{cls}"].append(cls_metrics["specificity"])

        metrics["balanced_accuracy"].append(balanced_accuracy)

        image_metrics.append({"longitude": long,
            "latitude": lat,
            "iou": avg_iou,
            **{f"precision_{cls}": per_class_metrics[cls]["precision"] if cls in per_class_metrics else None for cls in classes},
            **{f"recall_{cls}": per_class_metrics[cls]["recall"] if cls in per_class_metrics else None for cls in classes},
            **{f"f1_score_{cls}": per_class_metrics[cls]["f1_score"] if cls in per_class_metrics else None for cls in classes},
            **{f"accuracy_{cls}": per_class_metrics[cls]["accuracy"] if cls in per_class_metrics else None for cls in classes},
            **{f"specificity_{cls}": per_class_metrics[cls]["specificity"] if cls in per_class_metrics else None for cls in classes},
            "balanced_accuracy": balanced_accuracy})

        print(f"Metrics for image {long}, {lat}: IoU={avg_iou}")
        for cls in classes:
            if cls in per_class_metrics:
                print(f"  {cls} - Precision={per_class_metrics[cls]['precision']}, Recall={per_class_metrics[cls]['recall']}, F1 Score={per_class_metrics[cls]['f1_score']}, Accuracy={per_class_metrics[cls]['accuracy']}, Specificity={per_class_metrics[cls]['specificity']}")
        print(f"Balanced Accuracy={balanced_accuracy}")

    overall_metrics = {key: np.nanmean([value for value in values if value is not None]) if len(values) > 0 else None for key, values in metrics.items()}
    overall_metrics["longitude"] = "Overall"
    overall_metrics["latitude"] = "Metrics"

    print("Overall Metrics:")
    print("0: Parking lot, 1: Private, 2: Public")
    for key, value in overall_metrics.items():
        if key not in {"longitude", "latitude"}:
            if value is not None:
                print(f"Average {key.replace('_', ' ').title()}: {value:.2f}")
            else:
                print(f"Average {key.replace('_', ' ').title()}: None")

    fieldnames = ["longitude", "latitude", "iou", *[f"{metric}_{cls}" for cls in classes for metric in ["precision", "recall", "f1_score", "accuracy", "specificity"]], "balanced_accuracy"]
    
    with open(output_file, mode="w", newline="") as file:
        writer = csv.DictWriter(file, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(image_metrics)
        writer.writerow(overall_metrics)

    print(f"Metrics saved to {output_file}")

if __name__ == "__main__":
    main("classification_test_set")