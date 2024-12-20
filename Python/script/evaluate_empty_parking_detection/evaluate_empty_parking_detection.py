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
                    cv2.polylines(img, [box_points], isClosed=True, color=(255, 0, 0), thickness=2) #blue if on the road
                else:
                    print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is not on the road (possibly parked)")
                    detections_parking.append([x_center, y_center, width, height, angle_degrees, orientation])
                    cv2.polylines(img, [box_points], isClosed=True, color=(0, 0, 255), thickness=2) #red if parked

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

    return avg_width_meters, avg_length_meters, float(avg_width_pixels), float(avg_length_pixels)

def detect_empty_spots(cars, avg_spot_width, avg_spot_length, avg_width_pixels, avg_length_pixels, gap_threshold_meters=12, duplicate_threshold_meters=1, overlap_threshold_meters=1.25):
    """
    Detects empty spots in a group of parked cars in an image, based on the detected car bounding box centers.
    There are 4 cases: horizontal in a row, horizontal stacked, vertical in a column and vertical side by side.
    Duplicate spots and spots coinciding with cars identified by the Yolo model are removed.
    
    Params:
        cars (list): List of car bounding box centers with orientation horizontal or vertical
        avg_spot_width (float): Average width of a parking spot in meters
        avg_spot_length (float): Average length of a parking spot in meters
        avg_width_pixels (float): Average width of a parking spot in pixels
        avg_length_pixels (float): Average length of a parking spot in pixels
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
                    empty_spots.append((empty_x_center, empty_y_center, avg_width_pixels, avg_length_pixels, angle_degrees, alignment))
                    print(f"Empty parking spot at {empty_x_center}, {empty_y_center}")

    find_empty_spots(horizontal_cars_sorted_by_lat, 'horizontal', avg_spot_width, gap_threshold_meters) #Horizontal spots in a row
    find_empty_spots(vertical_cars_sorted_by_long, 'vertical', avg_spot_width, gap_threshold_meters)  #Vertical spots in columns

    empty_spots = sorted(empty_spots, key=lambda spot: (spot[0], spot[1]))
    unique_empty_spots = []

    for i, spot in enumerate(empty_spots):
        if i == 0:
            unique_empty_spots.append(spot)
        else:
            distance_to_prev = geodesic((spot[1], spot[0]), (empty_spots[i - 1][1], empty_spots[i - 1][0])).meters
            
            if distance_to_prev >= duplicate_threshold_meters or spot[5] != empty_spots[i - 1][5]:
                unique_empty_spots.append(spot)
            #else:
                #print(f"Removed spot at {spot[0]}, {spot[1]} due to proximity to {empty_spots[i - 1][0]}, {empty_spots[i - 1][1]}. Distance: {distance_to_prev:.2f}")


    filtered_empty_spots = []

    for empty_spot in unique_empty_spots:
        empty_x, empty_y = empty_spot[0], empty_spot[1] 
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
    
    for long, lat, width, length, rotation, alignment in empty_spots:
        x_center, y_center = convert_coordinates_to_bounding_box(long, lat, center_long, center_lat)
        
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
            filtered_empty_spots.append((long, lat, width, length, rotation, alignment))

    return filtered_empty_spots


def draw_empty_spots_on_image_original(image_path, empty_spots, center_long, center_lat, avg_spot_width, avg_spot_length):
    """
    Original function which seems to work better even though it doesn't take the rotation into account
    Draws the empty parking spots on the image

    Params:
        image_path (str): Path to the image
        empty_spots (list): List of empty parking spots' center coordinates
        center_long (float): Longitude of the center of the image
        center_lat (float): Latitude of the center of the image.
        avg_spot_width (float): Average width of a parking spot in pixels
        avg_spot_length (float): Average length of a parking spot in pixels
    """
    image = cv2.imread(image_path)

    for long, lat, _, _, _, alignment  in empty_spots:
        x_pixel, y_pixel = convert_coordinates_to_bounding_box(long, lat, center_long, center_lat)
        
        if alignment == 'horizontal':
            x1 = int(x_pixel - avg_spot_width // 2)
            y1 = int(y_pixel - avg_spot_length // 2)
            x2 = int(x_pixel + avg_spot_width // 2)
            y2 = int(y_pixel + avg_spot_length // 2)
        else: 
            x1 = int(x_pixel - avg_spot_length // 2)
            y1 = int(y_pixel - avg_spot_width // 2)
            x2 = int(x_pixel + avg_spot_length // 2)
            y2 = int(y_pixel + avg_spot_width // 2)  
        
        cv2.rectangle(image, (x1, y1), (x2, y2), (0, 255, 0), 2)

    cv2.imwrite(image_path, image)


def get_empty_parking_in_image(model, longitude, latitude, directory):
    """
    Detects all the parking spaces in the image (at longitude/latitude) and returns a list of coordinates, agregating
    the cars (of the road) found by the Yolo model and the empty parking spots found (which are drawn and added to the image).
    It returns only the empty parking spots found for evaluation. 

    Params:
        model : YOLO model
        longitude (float): Longitude value
        latitude (float): Latitude value
        directory(str): Path to the directory containing the images and the labels in json format 


    Returns: 
        all_detections (list): List of all empty parking spots detected in the image in the format log, lat, width, height, angle, orientation
    """
    output_folder = directory
    output_path_satelite_image = os.path.join(output_folder, f'{longitude}_{latitude}_satellite.png')
    output_path_road_image = os.path.join(output_folder, f'{longitude}_{latitude}_road.png')
    output_path_mask_image = os.path.join(output_folder, f'{longitude}_{latitude}_mask.png')
    output_path_bb_image = os.path.join(output_folder, f'{longitude}_{latitude}_bounding_boxes.png')

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
        empty_spots = detect_empty_spots(all_detections, avg_width_meters, avg_length_meters, avg_width_pixels, avg_length_pixels)
        empty_spots_filtered = filter_empty_spots_on_road(empty_spots, output_path_mask_image, longitude, latitude, avg_width_pixels, avg_length_pixels)
        #print(f'Filtered spots after road mask: {len(empty_spots_filtered)}')
        draw_empty_spots_on_image_original(output_path_bb_image, empty_spots_filtered, longitude, latitude, avg_width_pixels, avg_length_pixels)
    
    return empty_spots_filtered


def get_predictions_in_image(model, long, lat, directory):
    """
    Returns all empty parking spot predictions in an image in the correct format (with the pixel coordinates needed for evaluation)

    Params:
        model : YOLO model
        longitude (float): Longitude value of the image
        latitude (float): Latitude value of the image
        directory(str): Path to the directory containing the images and the labels in json format 


    Returns:
            empty_detections (list): List of all the empty parking spots found in the image in the format x, y, width, height, orientation (x and y being pixel values)
    """
    empty_detections = []

    detections = get_empty_parking_in_image(model, long, lat, directory)
    for x_center, y_center, width, height, _, orientation in detections:
        x_pixel, y_pixel =  convert_coordinates_to_bounding_box(x_center, y_center, long, lat)
        empty_detections.append([x_pixel, y_pixel, width, height, orientation])

    return empty_detections


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
                
                if len(parts) != 6:
                    print(f"Warning: Skipping line due to unexpected format: {line}")
                    continue
                
                try:
                    x_pixel = float(parts[1])*image_width #the values given by label studio are normalized and we want the denormalized values to compare with the predictions
                    y_pixel = float(parts[2])*image_height
                    width = float(parts[3])*image_width
                    height = float(parts[4])*image_height
                    orientation = parts[5]
                    true_labels.append([x_pixel, y_pixel, width, height, orientation])
                except ValueError:
                    print(f"Warning: Invalid data format in line: {line}")
                    continue
    except IOError as e:
        print(f"Error reading file {true_labels_file}: {e}")

    return true_labels

def get_true_labels_automatic_orientation_labelling(long, lat, directory, image_width=400, image_height=400):
    """
    Retrieve the true labels for a specific image.
    The orientation is automatically labelled and the writen in the file

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
    updated_lines = []

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
                    if width >= height:
                        orientation = "horizontal"
                    else:
                        orientation = "vertical"
                    true_labels.append([x_pixel, y_pixel, width, height, orientation])

                    updated_line = f"{parts[0]} {parts[1]} {parts[2]} {parts[3]} {parts[4]} {orientation}"
                    updated_lines.append(updated_line)

                except ValueError:
                    print(f"Warning: Invalid data format in line: {line}")
                    continue

        with open(file_path, 'w') as file:
            file.write("\n".join(updated_lines))
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

    for x_pixel, y_pixel, width, height, _ in true_labels:#here we don't need the orientation as in the labelling the width and height are correct (but in the model predictions width is almost always larger than height)
        
        x1 = int(x_pixel - width // 2)
        y1 = int(y_pixel - height // 2)
        x2 = int(x_pixel + width // 2)            
        y2 = int(y_pixel + height // 2)

        cv2.rectangle(image, (x1, y1), (x2, y2), (180, 105, 255), 2)

    cv2.imwrite(image_path, image)

def evaluate_predictions(predictions, true_labels, iou_threshold=0.35):
    """
    Evaluates our predictions in an image with the true labels
    Returns Average IoU, Precision, Recall, F1 score, Orientation Accuracy, 
    Spot Detection Ratio, Spot Detection Error, False Positive Rate, False Negative Rate

    Params:
        predictions (list): List of predictions
        true_labels (list): List of true_labels
        iou_threshold (float): IoU threshold to consider a prediction correct

    Returns:
        avg_iou, precision, recall, f1_score, orientation_accuracy, spot_detection_ratio, spot_detection_error, fpr, fnr (float): Average IoU, Precision, Recall, F1 score, Orientation Accuracy, Spot Detection Ratio, Spot Detection Error, FPR, FNR
    """

    def calculate_iou(box1, box2):
        """Calculates IoU for two bounding boxes.
        Box1 is the predictions and box2 is the true label"""
        x1, y1, w1, h1, orientation = box1
        x2, y2, w2, h2, _ = box2

        if orientation == 'vertical':
            w1, h1 = h1, w1  #swap width and height in the vertical case for the model predictions

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

    true_positives = 0
    false_positives = 0
    false_negatives = 0
    iou_scores = []
    orientation_matches = 0

    matched_labels = set()

    #For each predictions find the closest label and compare
    for pred in predictions:
        best_iou = 0
        best_match = None
        for i, gt in enumerate(true_labels):
            iou = calculate_iou(pred, gt)
            if iou > best_iou:
                best_iou = iou
                best_match = i

        if best_iou >= iou_threshold:
            if best_match not in matched_labels:
                true_positives += 1
                iou_scores.append(best_iou)
                matched_labels.add(best_match)

                if pred[4] == true_labels[best_match][4]:
                    orientation_matches += 1

    false_positives = len(predictions) - true_positives
    false_negatives = len(true_labels) - len(matched_labels)

    precision = true_positives / (true_positives + false_positives) if (true_positives + false_positives) > 0 else 0
    recall = true_positives / (true_positives + false_negatives) if (true_positives + false_negatives) > 0 else 0
    f1_score = 2 * (precision * recall) / (precision + recall) if (precision + recall) > 0 else 0
    avg_iou = sum(iou_scores) / len(iou_scores) if iou_scores else 0
    orientation_accuracy = orientation_matches / len(true_labels) if len(true_labels) > 0 else 0

    nb_predictions = len(predictions)
    nb_true_labels = len(true_labels)
    spot_detection_ratio = nb_predictions / nb_true_labels if nb_true_labels > 0 else 0
    spot_detection_error = abs(nb_predictions - nb_true_labels)

    fpr = false_positives / (false_positives + true_positives) if (false_positives + true_positives) > 0 else 0
    fnr = false_negatives / (false_negatives + true_positives) if (false_negatives + true_positives) > 0 else 0

    #print(true_positives, false_positives, false_negatives)

    #for labels in true_labels:
    #    for pred in predictions:
    #        iou = calculate_iou(labels, pred)
    #        print(f"True Labels: {labels}, Predictions: {pred}, IoU: {iou}")

    return avg_iou, precision, recall, f1_score, orientation_accuracy, spot_detection_ratio, spot_detection_error, fpr, fnr


def main(directory, output_file="metrics_empty_parking_detection.csv"):
    """
    Main function to evaluate the empty parking detction on the test images and calculate the corresponding performance metrics
    Saves Averages of IoU, Precision, Recall, F1 score, Orientation Accuracy, 
    Spot Detection Ratio, Spot Detection Error, False Positive Rate and False Negative Rate in a csv file for each
    image as well as the overall average metrics.

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
        "precision": [],
        "recall": [],
        "f1_score": [],
        "orientation_accuracy": [],
        "spot_detection_ratio": [],
        "spot_detection_error": [],
        "fpr": [],
        "fnr": []
    }

    image_metrics = []

    for long, lat in set(coordinates):
        predictions = get_predictions_in_image(model, long, lat, directory)
        true_labels = get_true_labels(long, lat, directory)
        #true_labels = get_true_labels_automatic_orientation_labelling(long, lat, directory)
        if not predictions and not true_labels:# when there are no detections and no true labels we want to skip the calculation of the metrics
            continue
        draw_true_labels(true_labels, directory, long, lat)
        #print(predictions)
        #print(true_labels)

        iou, precision, recall, f1_score, orientation_accuracy, spot_detection_ratio, spot_detection_error, fpr, fnr = evaluate_predictions(predictions, true_labels)

        metrics["iou"].append(iou)
        metrics["precision"].append(precision)
        metrics["recall"].append(recall)
        metrics["f1_score"].append(f1_score)
        metrics["orientation_accuracy"].append(orientation_accuracy)
        metrics["spot_detection_ratio"].append(spot_detection_ratio)
        metrics["spot_detection_error"].append(spot_detection_error)
        metrics["fpr"].append(fpr)
        metrics["fnr"].append(fnr)

        image_metrics.append({
            "longitude": long,
            "latitude": lat,
            "iou": iou,
            "precision": precision,
            "recall": recall,
            "f1_score": f1_score,
            "orientation_accuracy": orientation_accuracy,
            "spot_detection_ratio": spot_detection_ratio,
            "spot_detection_error": spot_detection_error,
            "fpr": fpr,
            "fnr": fnr
        })

        print(f"Metrics for image {long}, {lat}: IoU={iou:.2f}, Precision={precision:.2f}, Recall={recall:.2f}, F1 Score={f1_score:.2f}, Orientation Accuracy={orientation_accuracy:.2f}, Spot Detection Ratio={spot_detection_ratio:.2f}, Spot Detection Error={spot_detection_error:.2f}, FPR={fpr:.2f}, FNR={fnr:.2f}")

    avg_iou = sum(metrics["iou"]) / len(metrics["iou"])
    avg_precision = sum(metrics["precision"]) / len(metrics["precision"])
    avg_recall = sum(metrics["recall"]) / len(metrics["recall"])
    avg_f1 = sum(metrics["f1_score"]) / len(metrics["f1_score"])
    avg_orientation_accuracy = sum(metrics["orientation_accuracy"]) / len(metrics["orientation_accuracy"])
    avg_SDR = sum(metrics["spot_detection_ratio"]) / len(metrics["spot_detection_ratio"])
    avg_SDE = sum(metrics["spot_detection_error"]) / len(metrics["spot_detection_error"])
    avg_FPR = sum(metrics["fpr"]) / len(metrics["fpr"])
    avg_FNR = sum(metrics["fnr"]) / len(metrics["fnr"])

    overall_metrics = {
        "longitude": "Overall",
        "latitude": "Metrics",
        "iou": avg_iou,
        "precision": avg_precision,
        "recall": avg_recall,
        "f1_score": avg_f1,
        "orientation_accuracy": avg_orientation_accuracy,
        "spot_detection_ratio": avg_SDR,
        "spot_detection_error": avg_SDE,
        "fpr": avg_FPR,
        "fnr": avg_FNR
    }

    print("Overall Metrics:")
    print(f"Average IoU: {avg_iou:.2f}")
    print(f"Average Precision: {avg_precision:.2f}")
    print(f"Average Recall: {avg_recall:.2f}")
    print(f"Average F1 Score: {avg_f1:.2f}")
    print(f"Average Orientation Accuracy: {avg_orientation_accuracy:.2f}")
    print(f"Average Spot Detection Ratio (SDR): {avg_SDR:.2f}")
    print(f"Average Spot Detection Error (SDE): {avg_SDE:.2f}")
    print(f"Average False Positive Rate (FPR): {avg_FPR:.2f}")
    print(f"Average False Negative Rate (FNR): {avg_FNR:.2f}")

    with open(output_file, mode="w", newline="") as file:
        writer = csv.DictWriter(file, fieldnames=[
            "longitude", "latitude", "iou", "precision", "recall", "f1_score", "orientation_accuracy",
            "spot_detection_ratio", "spot_detection_error", "fpr", "fnr"
        ])
        writer.writeheader()
        writer.writerows(image_metrics)
        writer.writerow(overall_metrics)

    print(f"Metrics saved to {output_file}")

if __name__ == "__main__":
    main("all_test_images")