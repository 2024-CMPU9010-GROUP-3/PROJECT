import torch
from ultralytics import YOLO
import cv2
import requests
from PIL import Image
from io import BytesIO
import os
import numpy as np
import math
import pandas as pd
import random
from geopy.distance import geodesic
from sklearn.cluster import DBSCAN


def get_images(imag_save_path, longitude, latitude, mapbox_type):
    """
    Gets and saves an image from the Mapbox Static Images API

    Params:
        imag_save_path (str): Path to save the image
        longitude (float): Longitude value
        latitude (float): Latitude value
        mapbox_type (str): Satelite (satellite-v9) or Road view (streets-v12)
    """

    url = f'https://api.mapbox.com/styles/v1/mapbox/{mapbox_type}/static/{longitude},{latitude},18,0,0/400x400?access_token=pk.eyJ1Ijoia2F1c3R1Ymh0cml2ZWRpIiwiYSI6ImNtMWo2NndsbzB4N3EycHM1aGF2cDd5NzkifQ.4aegzX6Kfy3zW8pHkLWU7Q'
    response = requests.get(url)
    if response.status_code == 200:
        img = Image.open(BytesIO(response.content))
        img.save(imag_save_path)
        print(f"Image saved to {imag_save_path}")

def create_mask_using_canny(image_path, save_path, low_threshold=50, high_threshold=200):
   """
    Creates and saves a binary mask from the mapbox image of the road (Mapbox Streets) using Canny edge detction
    Gives worse results than the original and previous mask as there are too many edges interfering with the road

    Params:
        image_path (str): Path of the image
        save_path (str): Path to save the mask
        low_threshold (int): Lower threshold for Canny edge detection
        high_threshold (int): Upper threshold for Canny edge detection
    """
   img = cv2.imread(image_path)
   img_gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
   _, road_mask = cv2.threshold(img_gray, 240, 255, cv2.THRESH_BINARY)

   img_hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
   lower_orange = np.array([10, 100, 100])
   upper_orange = np.array([25, 255, 255])
   lower_yellow = np.array([25, 100, 100])
   upper_yellow = np.array([35, 255, 255])

   orange_mask = cv2.inRange(img_hsv, lower_orange, upper_orange)
   yellow_mask = cv2.inRange(img_hsv, lower_yellow, upper_yellow)
   combined_mask = cv2.bitwise_or(road_mask, orange_mask)
   combined_mask = cv2.bitwise_or(combined_mask, yellow_mask)

   masked_img = cv2.bitwise_and(img_gray, img_gray, mask=combined_mask)
   edges = cv2.Canny(masked_img, low_threshold, high_threshold)

   kernel = np.ones((3, 3), np.uint8)
   edges_cleaned = cv2.dilate(edges, kernel, iterations=1)
   edges_cleaned = cv2.erode(edges_cleaned, kernel, iterations=1)
   edges_cleaned = cv2.morphologyEx(edges_cleaned, cv2.MORPH_CLOSE, kernel)

   contours, _ = cv2.findContours(edges_cleaned, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
   mask_filtered = np.zeros_like(edges)

   for contour in contours:
        if cv2.contourArea(contour) > 50: 
            cv2.drawContours(mask_filtered, [contour], -1, 255, thickness=cv2.FILLED)

   cv2.imwrite(save_path, mask_filtered)


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
    combined_mask = cv2.bitwise_or(road_mask, orange_mask)
    combined_mask = cv2.bitwise_or(combined_mask, yellow_mask)

    kernel = np.ones((2, 2), np.uint8)#use smaller kernel as it works better
    combined_mask = cv2.morphologyEx(combined_mask, cv2.MORPH_CLOSE, kernel)#cv2.MORPH_CLOSE actually works better

    contours, _ = cv2.findContours(combined_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    mask_filtered = np.zeros_like(combined_mask)

    for contour in contours:
        if cv2.contourArea(contour) > 300: 
            cv2.drawContours(mask_filtered, [contour], -1, 255, thickness=cv2.FILLED)

    cv2.imwrite(save_path, mask_filtered)

def old_mask(image_path, save_path, threshold=240):
    """
    Creates and saves a binary mask from the mapbox image of the road (Mapbox Streets). The roads are in white while the rest of the image is darker
    Initial mask, that doesn't remove the street names

    Params:
        image_path (str): Path of the image
        save_path (str): Path to save the mask
        threshold (int): Threshold to differentiate the road from the areas outside of the road
    """
    img = cv2.imread(image_path)
    img_gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    _, road_mask = cv2.threshold(img_gray, threshold, 255, cv2.THRESH_BINARY)
    cv2.imwrite(save_path, road_mask)

def new_mask(image_path, save_path, threshold=240):
    """
    Creates and saves a binary mask from the mapbox image of the road (Mapbox Streets).
    The roads are in white and some additional roads are in orange/yellow (highways/ roads with more lanes).

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

    combined_mask = cv2.bitwise_or(road_mask, orange_mask)
    combined_mask = cv2.bitwise_or(combined_mask, yellow_mask)

    cv2.imwrite(save_path, combined_mask)

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
    
    horizontal_cars_sorted_by_long = sorted([car for car in cars if car[5] == 'horizontal'], key=lambda point: point[0]) 
    horizontal_cars_sorted_by_lat = sorted([car for car in cars if car[5] == 'horizontal'], key=lambda point: point[1])
    vertical_cars_sorted_by_long = sorted([car for car in cars if car[5] == 'vertical'], key=lambda point: point[0])  
    vertical_cars_sorted_by_lat = sorted([car for car in cars if car[5] == 'vertical'], key=lambda point: point[1]) 

    def find_empty_spots(sorted_cars, alignment, gap_dimension, gap_threshold_meters):
        """ Detects empty spots in the sorted list of cars for a specific alignment (horizontal or vertical) 
        
        Params:
        sorted_cars (list): List of car bounding box centers sorted by orientation (horizontal/vertical) and coordinates (long/lat)
        alignment (str): Either horizontal or vertical
        gap_dimension (str): Either avg_spot_width or avg_spot_length
        gap_threshold_meters (float): Maximum allowed gap to consider there is an empty parking spot or multiple parking spots
        """
        for i in range(len(sorted_cars) - 1):
            x_current, y_current, _, _, _, _ = sorted_cars[i]
            x_next, y_next, _, _, _, _ = sorted_cars[i + 1]

            gap_distance = geodesic((y_current, x_current), (y_next, x_next)).meters
            avg_half_dim = gap_dimension / 2
            adjusted_gap = gap_distance - 2 * avg_half_dim

            angle_radians = math.atan2(y_next - y_current, x_next - x_current)
            angle_degrees = math.degrees(angle_radians)
            
            if adjusted_gap <= gap_threshold_meters and adjusted_gap > gap_dimension:
                num_spots = int(adjusted_gap // gap_dimension)
                
                for j in range(1, num_spots + 1):
                    empty_x_center = x_current + j * (x_next - x_current) / (num_spots + 1)
                    empty_y_center = y_current + j * (y_next - y_current) / (num_spots + 1)
                    empty_spots.append(([empty_x_center, empty_y_center], angle_degrees, alignment))
                    print(f"Empty parking spot at {empty_x_center}, {empty_y_center}")

    find_empty_spots(horizontal_cars_sorted_by_long, 'horizontal', avg_spot_length, gap_threshold_meters) #Horizontal spots in a row
    find_empty_spots(horizontal_cars_sorted_by_lat, 'horizontal', avg_spot_width, gap_threshold_meters=9 ) #Horizontal spots stacked in a column
    find_empty_spots(vertical_cars_sorted_by_lat, 'vertical', avg_spot_length, gap_threshold_meters) #Vertical spots in columns
    find_empty_spots(vertical_cars_sorted_by_long, 'vertical', avg_spot_width, gap_threshold_meters=9)  # Vertical spots side by side in a row

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
            filtered_empty_spots.append((spot, rotation, alignment))

    return filtered_empty_spots

def draw_empty_spots_on_image(image_path, empty_spots, center_long, center_lat, avg_spot_width, avg_spot_length):
    """
    Draws the empty parking spots on the image

    Params:
        image_path (str): Path to the image
        empty_spots (list): List of empty parking spots' center coordinates, roation and alignment
        center_long (float): Longitude of the center of the image
        center_lat (float): Latitude of the center of the image
        avg_spot_width (float): Average width of a parking spot in pixels
        avg_spot_length (float): Average length of a parking spot in pixels
    """
    
    image = cv2.imread(image_path)

    for i, (spot, rotation, orientation) in enumerate(empty_spots):
        x_pixel, y_pixel = convert_coordinates_to_bounding_box(spot[0], spot[1], center_long, center_lat)

        width = avg_spot_width if orientation == 'horizontal' else avg_spot_length
        height = avg_spot_length if orientation == 'horizontal' else avg_spot_width
        
        if i > 0:
            prev_rotation = empty_spots[i - 1][1]
            rotation_diff = abs(rotation - prev_rotation)

            if rotation_diff < 15:
                rotation = prev_rotation
            elif abs(rotation % 90) < 10:
                rotation = round(rotation / 90) * 90

        rect = ((x_pixel, y_pixel), (width, height), rotation)
        box_points = cv2.boxPoints(rect)
        box_points = np.int32(box_points)

        cv2.polylines(image, [box_points], isClosed=True, color=(0, 255, 0), thickness=2)

    cv2.imwrite(image_path, image)


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

    for spot, _, orientation in empty_spots:
        x_pixel, y_pixel = convert_coordinates_to_bounding_box(spot[0], spot[1], center_long, center_lat)
        
        if orientation == 'horizontal':
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
    cluster_labels = []

    road_mask = cv2.imread(road_mask_path, cv2.IMREAD_GRAYSCALE)
    road_contours, _ = cv2.findContours(road_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    #print(center_long, center_lat)
    pixel_coords = []

    for spot in all_parking_spots:
        x_center, y_center = convert_coordinates_to_bounding_box(spot[0], spot[1], center_long, center_lat)
        pixel_coords.append([x_center, y_center])

    clustering = DBSCAN(eps=clustering_eps, min_samples=clustering_min_samples).fit(pixel_coords)
    labels = clustering.labels_

    for idx, spot in enumerate(all_parking_spots):
        x_center, y_center = pixel_coords[idx]
        cluster_label = labels[idx]
        cluster_labels.append(cluster_label)

        min_distance = float("inf")

        for contour in road_contours:
            distance = cv2.pointPolygonTest(contour, (x_center, y_center), measureDist=True)
            min_distance = min(min_distance, abs(distance))

        classification = "private"

        if min_distance <= road_proximity_threshold:
            classification = "public"

        if cluster_label != -1:
            cluster_size = np.sum(labels == cluster_label)
            if cluster_size >= parking_lot_min_spots:
                classification = "parking lot"

        classified_spots.append([spot[0], spot[1], classification])
        #print(classification)

    return classified_spots, cluster_labels

def draw_clusters_and_labels(image_path, spots, cluster_labels, center_long, center_lat):
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

    unique_clusters = set(cluster_labels)
    cluster_colors = {cluster: tuple(random.randint(0, 255) for _ in range(3)) for cluster in unique_clusters}

    classification_colors = {
        "public": (0, 255, 0),  #green
        "private": (0, 0, 255), #red
        "parking_lot": (255, 0, 0) #blue
    }

    for i, spot in enumerate(spots):
        cluster_label = cluster_labels[i]
        classification = spot[2]

        x, y = convert_coordinates_to_bounding_box(spot[0], spot[1], center_long, center_lat)
        x, y = int(x), int(y)

        cluster_color = cluster_colors.get(cluster_label, (255, 255, 255))
        classification_color = classification_colors.get(classification, (255, 255, 255))

        cv2.circle(image, (x, y), 5, cluster_color, -1)
        label = f"{classification}"
        cv2.putText(image, label, (x + 10, y - 10), cv2.FONT_HERSHEY_SIMPLEX, 0.6, classification_color, 2)

    cv2.imwrite(image_path, image)

def get_parking_coords_in_image(model, longitude, latitude):
    """
    Detects all the parking spaces in the image (at longitude/latitude) and returns a list of coordinates, agregating
    the cars (of the road) found by the Yolo model and the empty parking spots found (which are drawn and added to the image) 

    Params:
        model : YOLO model
        longitude (float): Longitude value
        latitude (float): Latitude value

    Returns: 
        all_detections (list): List of all coordinates of parking spots found in the image in the format long, lat, classification
    """
    output_folder = 'image_output'
    output_path_satelite_image = os.path.join(output_folder, f'{longitude}_{latitude}_satelite.png')
    output_path_road_image = os.path.join(output_folder, f'{longitude}_{latitude}_road.png')
    output_path_mask_image = os.path.join(output_folder, f'{longitude}_{latitude}_mask.png')
    output_path_bb_image = os.path.join(output_folder, f'{longitude}_{latitude}_bounding_boxes.png')

    get_images(output_path_satelite_image, longitude, latitude, 'satellite-v9')
    get_images(output_path_road_image, longitude, latitude, 'streets-v12')

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
        draw_empty_spots_on_image_original(output_path_bb_image, empty_spots_filtered, longitude, latitude, avg_width_pixels, avg_length_pixels)
        empty_spots_coords = [spot for spot, _, _ in empty_spots_filtered]
        all_detections.extend(empty_spots_coords)
        all_detections, cluster_labels = classify_parking_spots(all_detections, output_path_mask_image, longitude, latitude)
        draw_clusters_and_labels(output_path_bb_image, all_detections, cluster_labels, longitude, latitude)

    return all_detections

def long_lat_to_tile_coords(long, lat):
    """
    Converts longitude and latitude to tile coordinates

    Params:
        long (float): Longitude value
        lat (float): Latitude value

    Returns: 
        x_tile, y_tile (float): Tile coordinates
    """
    num_tiles = 2 ** 18
    lat_radians = math.radians(lat)
    x_tile = (long + 180.0) / 360.0 * num_tiles
    y_tile = (1.0 - math.log(math.tan(lat_radians) + (1 / math.cos(lat_radians))) / math.pi) / 2.0 * num_tiles
    
    return x_tile, y_tile

def tile_coords_to_long_lat(x_tile, y_tile):
    """
    Converts tile coordinates to longitude and latitude

    Params:
        x_tile (float): x-coordinate of the tile
        y_tile (float): y-coordinate of the tile

    Returns: 
        long, lat (float): Longitude and Latitude values
    """
    num_tiles = 2 ** 18
    long = x_tile / num_tiles * 360.0 - 180.0
    lat_rad = math.atan(math.sinh(math.pi * (1 - 2 * y_tile / num_tiles)))
    lat = math.degrees(lat_rad)

    return long, lat

    
def get_image_center_coords_from_bb(top_left_longitude, top_left_latitude, bottom_right_longitude, bottom_right_latitude):
    """
    Returns all the centers coordinates of each image necessary to generate within a bounding box

    Params:
        top_left_longitude (float): Longitude of the top left corner of the bounding box
        top_left_latitude (float): Latitude of the top left corner of the bounding box
        bottom_right_longitude (float): Longitude of the bottom right corner of the bounding box
        bottom_right_latitude (float): Latitude of the bottom right corner of the bounding box

    Returns: 
        centers (list):  List of all the center coordinates
    """
    top_left_x_tile, top_left_y_tile = long_lat_to_tile_coords(top_left_longitude, top_left_latitude)
    bottom_right_x_tile, bottom_right_y_tile = long_lat_to_tile_coords(bottom_right_longitude, bottom_right_latitude)

    num_hor_tiles = math.ceil(abs(bottom_right_x_tile - top_left_x_tile))
    num_vert_tiles = math.ceil(abs(top_left_y_tile - bottom_right_y_tile))

    print(f"Top-left tile coords: ({top_left_x_tile}, {top_left_y_tile})")
    print(f"Bottom-right tile coords: ({bottom_right_x_tile}, {bottom_right_y_tile})")
    print(f"Number of horizontal tiles: {num_hor_tiles}")
    print(f"Number of vertical tiles: {num_vert_tiles}")

    centers = []

    for i in range(num_hor_tiles):
        center_x_tile = top_left_x_tile + i + 0.5  #even in negative cases, we add (just means we get closer to 0)
        
        for j in range(num_vert_tiles):
            if top_left_y_tile > bottom_right_y_tile:
                center_y_tile = top_left_y_tile - j - 0.5 #in the normal case we decrease
            else:
                center_y_tile = top_left_y_tile + j + 0.5 #but if the latitude is negative we need to increase

            center_long, center_lat = tile_coords_to_long_lat(center_x_tile, center_y_tile)
            centers.append((center_long, center_lat))

    return centers


def main(top_left_longitude, top_left_latitude, bottom_right_longitude, bottom_right_latitude):
    """
    Detects all the parking spaces contained within a bounding box defined by the top left/ bottom right longitudes and latitudes.
    And saves them in a csv file.

    Params:
        top_left_longitude (float): Longitude of the top left corner of the bounding box
        top_left_latitude (float): Latitude of the top left corner of the bounding box
        bottom_right_longitude (float): Longitude of the bottom right corner of the bounding box
        bottom_right_latitude (float): Latitude of the bottom right corner of the bounding box

    """
    if not os.path.exists("image_output"):
        os.makedirs("image_output")

    model = YOLO("best - obb.pt")

    centers = get_image_center_coords_from_bb(top_left_longitude, top_left_latitude, bottom_right_longitude, bottom_right_latitude)

    all_detections = []

    for long, lat in centers:
        detections = get_parking_coords_in_image(model, long, lat)
        for detection in detections:
            all_detections.append([detection[0], detection[1], detection[2]])

    df = pd.DataFrame(all_detections, columns=["longitude", "latitude", "type"])
    df = df.drop_duplicates(subset=["longitude", "latitude"], keep="first")# remove duplicate coords as there is potential overlap in the images
    df.to_csv(f"coordinates_in_{top_left_longitude}_{top_left_latitude}-{bottom_right_longitude}_{bottom_right_latitude}.csv", index=False)
    

if __name__ == "__main__":
    #main(-6.2264, 53.4194, -6.2219, 53.4221)#parking lot
    #main(-6.2563, 53.3952, -6.2525, 53.3974)#residential area
    #main(-6.289, 53.3653, -6.2842, 53.3681)#residential area
    #main(-6.2737, 53.3436, -6.2709, 53.3452)#urban area
    #main(-6.2751, 53.347, -6.272, 53.3489)#urban area
    #main(-6.2844, 53.3589, -6.2816, 53.3606)#residential area
    #main(-6.2901, 53.3587, -6.2872, 53.3604)#residential area
    #main(-6.2859, 53.3636, -6.2823, 53.3656)#residential area
    #main(-6.2754, 53.3471, -6.2732, 53.3483)#urban area
    #main(-6.2652, 53.3525, -6.2625, 53.3541)#urban with parking lot

    main()

