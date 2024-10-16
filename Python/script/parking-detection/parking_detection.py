import torch
import cv2
import numpy as np

def create_mask(image_path, save_path, threshold=240):
    """
    Creates a binary mask from the mapbox image of the road (Mapbox Streets). The roads are in white while the rest of the image is darker.
    
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

def main():
    model = torch.hub.load('ultralytics/yolov5', 'custom', path='yolov5/runs/train/exp3/weights/best.pt', force_reload=True)
    create_mask('street3.png', 'mask.png')
    detections = detect_parking_spots_in_image('sateliteview3.png','mask.png', 'bounding_boxes_parked_cars.png', model)
    #print(detections)


if __name__ == "__main__":
    main()
