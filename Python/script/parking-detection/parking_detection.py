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


def create_mask(image_path, save_path, threshold=240):
    """
    Creates and saves a binary mask from the mapbox image of the road (Mapbox Streets). The roads are in white while the rest of the image is darker
    
    Params:
        image_path (str): Path of the image
        save_path (str): Path to save the mask
        threshold (int): Threshold to differentiate the road from the areas outside of the road
    """
    img = cv2.imread(image_path)
    img_gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    _, road_mask = cv2.threshold(img_gray, threshold, 255, cv2.THRESH_BINARY)
    kernel = np.ones((5, 5), np.uint8)
    road_mask_dilated = cv2.dilate(road_mask, kernel, iterations=2)
    
    contours, _ = cv2.findContours(road_mask_dilated, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    mask_filtered = np.zeros_like(road_mask)
    
    for contour in contours:
        if cv2.contourArea(contour) > 600:  
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

def detect_parking_spots_in_image(image_path, road_mask_path, output_image_path, model):
    """
    Detect cars in the image using the retrained YOLO model and remove those on the road

    Params:
        image_path (str): Path of the image
        road_mask_path (str): Path of the saved mask
        output_image_path (str): Path to save the image with bounding boxes, red for parking and blue cars on the road
        model_path : YOLO model.
        
    Returns:
        A lsit of the bounding boxes for the cars not on the road
    """
    img = cv2.imread(image_path)
    results = model.predict(img)

    detections_parking = []

    road_mask = cv2.imread(road_mask_path, cv2.IMREAD_GRAYSCALE)

    if not results:
        print("No cars detected.")
        return detections_parking

    for result in results:
        detections = result.boxes

        for box in detections:
            x_min, y_min, x_max, y_max = box.xyxy[0] 
            conf = box.conf[0] 
            cls = int(box.cls[0])

            if cls == 0:
                car_region = road_mask[int(y_min):int(y_max), int(x_min):int(x_max)]

                if car_region.size == 0:
                    continue

                road_pixels = cv2.countNonZero(car_region)
                total_pixels = car_region.size

                if road_pixels / total_pixels > 0.5:
                    print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is on the road")
                    cv2.rectangle(img, (int(x_min), int(y_min)), (int(x_max), int(y_max)), (255, 0, 0), 2)  #blue if on the road
                else:
                    print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is not on the road (possibly parked)")
                    detections_parking.append([x_min, y_min, x_max, y_max])
                    cv2.rectangle(img, (int(x_min), int(y_min)), (int(x_max), int(y_max)), (0, 0, 255), 2)  #red if parked

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
    

def get_center_bounding_box(x_min, y_min, x_max, y_max):
    """
    Returns the center of the bounding box
    
    Params:
        x_min, y_min, x_max, y_max (int): Top left and bottom right cordinates of the bounding box
    """
    x = (x_min + x_max) / 2
    y = (y_min + y_max) / 2

    return x, y

def get_parking_coords_in_image(model, longitude, latitude):
    """
    Detects the parking spaces in the image (at longitude/latitude) and returns a list of coordinates

    Params:
        longitude (float): Longitude value
        latitude (float): Latitude value
    """
    output_folder = 'image_output'
    output_path_satelite_image = os.path.join(output_folder, f'{longitude}_{latitude}_satelite.png')
    output_path_road_image = os.path.join(output_folder, f'{longitude}_{latitude}_road.png')
    output_path_mask_image = os.path.join(output_folder, f'{longitude}_{latitude}_mask.png')
    output_path_bb_image = os.path.join(output_folder, f'{longitude}_{latitude}_bounding_boxes.png')

    get_images(output_path_satelite_image, longitude, latitude, 'satellite-v9')
    get_images(output_path_road_image, longitude, latitude, 'streets-v12')

    old_mask(output_path_road_image, output_path_mask_image)
    detections = detect_parking_spots_in_image(output_path_satelite_image, output_path_mask_image, output_path_bb_image, model)

    all_detections = []

    for detection in detections:
        x, y = get_center_bounding_box(detection[0], detection[1], detection[2], detection[3])
        long, lat = convert_bounding_box_to_coordinates(x, y, longitude, latitude)
        print(f"Car coordinates: ({long}, {lat})")
        all_detections.append([long, lat])

    return all_detections

def long_lat_to_tile_coords(long, lat):
    num_tiles = 2 ** 18
    lat_radians = math.radians(lat)
    x_tile = (long + 180.0) / 360.0 * num_tiles
    y_tile = (1.0 - math.log(math.tan(lat_radians) + (1 / math.cos(lat_radians))) / math.pi) / 2.0 * num_tiles
    
    return x_tile, y_tile

def tile_coords_to_long_lat(x_tile, y_tile):
    num_tiles = 2 ** 18
    long = x_tile / num_tiles * 360.0 - 180.0
    lat_rad = math.atan(math.sinh(math.pi * (1 - 2 * y_tile / num_tiles)))
    lat = math.degrees(lat_rad)

    return long, lat

    
    return longitude, latitude

def get_image_center_coords_from_bb(top_left_longitude, top_left_latitude, bottom_right_longitude, bottom_right_latitude):
    """
    Returns centers of all the images within the bounding box.
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
        for j in range(num_vert_tiles):
            center_x_tile = top_left_x_tile + i + 0.5 
            center_y_tile = top_left_y_tile - j - 0.5
            center_long, center_lat = tile_coords_to_long_lat(center_x_tile, center_y_tile)
            centers.append((center_long, center_lat))

    return centers


def main(top_left_longitude, top_left_latitude, bottom_right_longitude, bottom_right_latitude):
    """
    Detects the parking spaces in the image at the longitude and latitude

    Params:
        longitude (float): Longitude value
        latitude (float): Latitude value
    """
    if not os.path.exists('image_output'):
        os.makedirs('image_output')

    #model = YOLO('best.pt')

    centers = get_image_center_coords_from_bb(top_left_longitude, top_left_latitude, bottom_right_longitude, bottom_right_latitude)
    """all_detections = []

    for long, lat in centers:
        detections = get_parking_coords_in_image(model, long, lat)
        for detection in detections:
            all_detections.append(detection) 

    df = pd.DataFrame(all_detections, columns=["longitude", "latitude"])
    df.to_csv(f"coordinates_in_{top_left_longitude}_{top_left_latitude}-{bottom_right_longitude}_{bottom_right_latitude}.csv", index=False)
    """
    print(centers)
if __name__ == "__main__":
    main(-6.2747, 53.3694, -6.2671, 53.3649)