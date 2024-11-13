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
from geopy.distance import geodesic


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

def detect_parking_spots_in_image(image_path, road_mask_path, output_image_path, model):
    """
    Detect cars in the image using the retrained YOLO model and remove those on the road

    Params:
        image_path (str): Path of the image
        road_mask_path (str): Path of the saved mask
        output_image_path (str): Path to save the image with bounding boxes, red for parking and blue cars on the road
        model : YOLO model.
        
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
                    print(orientation)
                    cv2.polylines(img, [box_points], isClosed=True, color=(255, 0, 0), thickness=2) #blue if on the road
                else:
                    print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is not on the road (possibly parked)")
                    print(orientation)
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

def get_center_bounding_box(x_min, y_min, x_max, y_max):
    """
    Returns the center of the bounding box
    
    Params:
        x_min, y_min, x_max, y_max (int): Top left and bottom right cordinates of the bounding box

    Returns:
        x, y (float): Center coordinates of the bounding box
    """
    x = (x_min + x_max) / 2
    y = (y_min + y_max) / 2

    return x, y

def calculate_avg_spot_dimensions(cars):
    """
    Calculates average parking spot width and length from the bounding box found as it is very variable
    
    Params:
        cars (list): List of corrdinates with width and length
        
    Returns:
        avg_width_meters, avg_length_meters, avg_width_pixels, avg_length_pixels (float): Average width and length of the cars identified in meters and in pixels
    """
    widths = [car[2] for car in cars]
    lengths = [car[3] for car in cars]
    avg_width_pixels = np.median(widths)
    avg_length_pixels = np.median(lengths)

    #cross product based on optimal values found previously
    avg_width_meters = 2 * avg_width_pixels / 18
    avg_length_meters = 3.2 * avg_length_pixels / 18

    print(avg_width_meters, avg_length_meters, avg_width_pixels, avg_length_pixels)
    return avg_width_meters, avg_length_meters, avg_width_pixels, avg_length_pixels

def detect_empty_spots(cars, avg_spot_width, avg_spot_length, gap_threshold_meters=12, duplicate_threshold_meters=1):
    """
    Detects empty spots in rows of parked cars based on detected car bounding box centers
    
    Params:
        cars (list): List of car bounding box centers with orientation horizontal or vertical
        avg_spot_width (float): Average width of a parking spot in meters
        avg_spot_length (float): Average length of a parking spot in meters
        gap_threshold_meters (float): Maximum allowed gap to consider there is an empty parking spot or multiple parking spots
        duplicate_threshold_meters (float): Threshold to differenciate between spots that are considered identical in meters

    Returns:
       empty_spots (list): List of coordinates of estimated empty parking spots with horizontal or vertical orientation (for drawing the boxes)
    """
    empty_spots = []
    
    horizontal_cars_sorted_by_long = sorted([car for car in cars if car[5] == 'horizontal'], key=lambda point: point[0]) 
    horizontal_cars_sorted_by_lat = sorted([car for car in cars if car[5] == 'horizontal'], key=lambda point: point[1])
    vertical_cars_sorted_by_long = sorted([car for car in cars if car[5] == 'vertical'], key=lambda point: point[0])  
    vertical_cars_sorted_by_lat = sorted([car for car in cars if car[5] == 'vertical'], key=lambda point: point[1]) 

    def find_empty_spots(sorted_cars, alignment, gap_dimension):
        """ Detects empty spots in the sorted list of cars for a specific alignment (horizontal or vertical) 
        
        Params:
        sorted_cars (list): List of car bounding box centers sorted by orientation (horizontal/vertical) and coordinates (long/lat)
        alignment (str): Either horizontal or vertical
        gap_dimension (str): Either avg_spot_width or avg_spot_length
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
                    print(f"Empty parking spot at {empty_x_center}, {empty_y_center}, {alignment}, {angle_degrees}")

    find_empty_spots(horizontal_cars_sorted_by_long, 'horizontal', avg_spot_length) #Horizontal spots in a row
    find_empty_spots(horizontal_cars_sorted_by_lat, 'horizontal', avg_spot_width) #Horizontal spots stacked in a column
    find_empty_spots(vertical_cars_sorted_by_long, 'vertical', avg_spot_length) #Vertical spots in columns
    find_empty_spots(vertical_cars_sorted_by_lat, 'vertical', avg_spot_width) # Vertical spots side by side in a row

    empty_spots = sorted(empty_spots, key=lambda spot: (spot[0][0], spot[0][1]))
    unique_empty_spots = []

    for i, spot in enumerate(empty_spots):
        if i == 0:
            unique_empty_spots.append(spot)
        else:
            distance_to_prev = geodesic((spot[0][1], spot[0][0]), (empty_spots[i - 1][0][1], empty_spots[i - 1][0][0])).meters
            
            if distance_to_prev >= duplicate_threshold_meters or spot[2] != empty_spots[i - 1][2]:
                unique_empty_spots.append(spot)

    return unique_empty_spots


def draw_empty_spots_on_image(image_path, empty_spots, center_long, center_lat, avg_spot_width, avg_spot_length):
    """
    Draws the empty parking spots on the image

    Params:
        image_path (str): Path to the image
        empty_spots (list): List of empty parking spots' center coordinates, roation and alignment
        center_long (float): Longitude of the center of the image
        center_lat (float): Latitude of the center of the image.
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


def get_parking_coords_in_image(model, longitude, latitude):
    """
    Detects the parking spaces in the image (at longitude/latitude) and returns a list of coordinates

    Params:
        model : YOLO model
        longitude (float): Longitude value
        latitude (float): Latitude value

    Returns: 
        all_detections (list): List of all coordinates of parking spots found in the image in the format log, lat, width, height, angle
    """
    output_folder = 'image_output'
    output_path_satelite_image = os.path.join(output_folder, f'{longitude}_{latitude}_satelite.png')
    output_path_road_image = os.path.join(output_folder, f'{longitude}_{latitude}_road.png')
    output_path_mask_image = os.path.join(output_folder, f'{longitude}_{latitude}_mask.png')
    output_path_bb_image = os.path.join(output_folder, f'{longitude}_{latitude}_bounding_boxes.png')

    #get_images(output_path_satelite_image, longitude, latitude, 'satellite-v9')
    #get_images(output_path_road_image, longitude, latitude, 'streets-v12')

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
        draw_empty_spots_on_image(output_path_bb_image, empty_spots, longitude, latitude, avg_width_pixels, avg_length_pixels)
        empty_spots_coords = [spot for spot, _, _ in empty_spots]
        all_detections.extend(empty_spots_coords)

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
            all_detections.append([detection[0], detection[1]])

    df = pd.DataFrame(all_detections, columns=["longitude", "latitude"])
    df = df.drop_duplicates(subset=["longitude", "latitude"], keep="first")# remove duplicate coords as there is potential overlap in the images
    df.to_csv(f"coordinates_in_{top_left_longitude}_{top_left_latitude}-{bottom_right_longitude}_{bottom_right_latitude}.csv", index=False)
    

if __name__ == "__main__":
    main(-6.2576, 53.3388, -6.2566, 53.3394)
    main(-6.2608, 53.3464, -6.2598, 53.347)
    main(-6.2617, 53.3462, -6.2606, 53.3469)
    main(-6.2854, 53.3511, -6.2843, 53.3517)
    main(-6.2893, 53.3486, -6.2883, 53.3492)
    main(-6.2899, 53.3473, -6.2889, 53.3479)
    main(-6.2903, 53.349, -6.2893, 53.3496) 
    main(-6.2657, 53.3567, -6.2646, 53.3574)

    main(-6.2853, 53.353, -6.2845, 53.3533)#hor
    main(-6.2847, 53.3526, -6.2839, 53.3529)
    main(-6.2845, 53.3524, -6.2837, 53.3527)
    main(-6.2821, 53.3492, -6.2814, 53.3495)#vert
    main(-6.2735, 53.3473, -6.2727, 53.3476)
    main(-6.2723, 53.3478, -6.2719, 53.3479 )#parking
    main(-6.2548, 53.3274, -6.2541, 53.3277)
    main(-6.2705, 53.3466, -6.2698, 53.3468)#vert
    main(-6.2646, 53.3432, -6.2641, 53.3434)
    main(-6.2649, 53.3385, -6.2641, 53.3388)
    main(-6.2611, 53.3308, -6.2604, 53.3311)
    main(-6.2512, 53.3252, -6.2506, 53.3254)
    main(-6.2708, 53.3456, -6.27, 53.3459)#hor
    main(-6.2653, 53.3392, -6.2645, 53.3395)#mix
    main(-6.2461, 53.3198, -6.2453, 53.3201)
    main(-6.2463, 53.3195, -6.2455, 53.3198)
    main(-6.2465, 53.3193, -6.2457, 53.3196)
    