from ultralytics import YOLO

# import model
yolo8_model = YOLO("yolov8n.yaml")

# train on dataset
yolo8_model.train(data="data.yaml", epochs=1)

# validate
yolo8_model.val(data="data.yaml")

