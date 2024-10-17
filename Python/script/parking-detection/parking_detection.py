import torch
import cv2
import requests
from PIL import Image
from io import BytesIO
import os
import numpy as np

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
    Creates and saves a binary mask from the mapbox image of the road (Mapbox Streets). The roads are in white while the rest of the image is darker.
    
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

def detect_parking_spots_in_image(image_path, road_mask_path, output_image_path, model):
    """
    Detect cars in the image using the retrained YOLO model and remove those on the road

    Parameters:
        image_path (str): Path of the image
        road_mask_path (str): Path of the saved mask
        output_image_path (str): Path to save the image with bounding boxes, red for parking and blue cars on the road
        model_path : YOLO model.
        
    Returns:
        a lsit of the bounding boxes for the cars not on the road
    """

    img = cv2.imread(image_path)

    results = model(img)
    detections = results.xyxy[0]  # bounding boxes format [x_min, y_min, x_max, y_max, confidence, class]
    detections_parking = []
    #print(detections)

    road_mask = cv2.imread(road_mask_path, cv2.IMREAD_GRAYSCALE)

    if detections.shape[0] == 0:
        print("No cars detected.")
        return detections_parking
    
    for detection in detections:
        x_min = int(detection[0].item())
        y_min = int(detection[1].item())
        x_max = int(detection[2].item())
        y_max = int(detection[3].item())
        conf = detection[4].item()
        cls = int(detection[5].item())
        #print(detection)

        if cls == 0:
            car_region = road_mask[y_min:y_max, x_min:x_max]

            if car_region.size == 0:
                continue

            road_pixels = cv2.countNonZero(car_region)
            total_pixels = car_region.size

            if road_pixels / total_pixels > 0.5:
                print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is on the road")
                cv2.rectangle(img, (x_min, y_min), (x_max, y_max), (255, 0, 0), 2)#blue if on the road
            else:
                print(f"Car at [{x_min}, {y_min}, {x_max}, {y_max}] is not on the road (possibly parked)")
                detections_parking.append([x_min, y_min, x_max, y_max, conf, cls])
                cv2.rectangle(img, (x_min, y_min), (x_max, y_max), (0, 0, 255), 2)#red if parked 

    cv2.imwrite(output_image_path, img)
    return detections_parking

def main(longitude, latidue):
    """
    Detects the parking spaces in the image at the longitude and latitude

    Params:
        longitude (float): Longitude value
        latitude (float): Latitude value
    """
    model = torch.hub.load('ultralytics/yolov5', 'custom', path='yolov5/runs/train/exp3/weights/best.pt', force_reload=True)
    
    if not os.path.exists('image_output'):
        os.makedirs('image_output')

    output_folder = 'image_output'
    output_path_satelite_image = os.path.join(output_folder, f'{longitude}_{latidue}_satelite.png')
    output_path_road_image = os.path.join(output_folder, f'{longitude}_{latidue}_road.png')
    output_path_mask_image = os.path.join(output_folder, f'{longitude}_{latidue}_mask.png')
    output_path_bb_image = os.path.join(output_folder, f'{longitude}_{latidue}_bounding_boxes.png')


    get_images(output_path_satelite_image, longitude, latidue, 'satellite-v9')
    get_images(output_path_road_image, longitude, latidue, 'streets-v12')

    create_mask(output_path_road_image, output_path_mask_image)
    detections = detect_parking_spots_in_image(output_path_satelite_image, output_path_mask_image, output_path_bb_image, model)
    #print(detections)

if __name__ == "__main__":
    main(-6.2849, 53.3531)