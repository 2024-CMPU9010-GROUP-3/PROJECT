### Instructions
I did this in VS code, the code is to be run in the terminal 

1. Clone the Yolov5 GitHub repo
git clone https://github.com/ultralytics/yolov5.git  

2. Cd into the cloned directory & fulfill requirements to run the project
cd yolov5
pip install -r requirements.txt

3. Download pre-trained weights for the yolov5 model (i've chosen the small model - yolov5s.pt -- to start with)
python detect.py --weights yolov5s.pt --source data/images

4. fine tune model by training & validating on our images
python train.py --img 640 --batch 16 --epochs 100 --data data.yaml --weights yolov5s.pt --cache

5. view results of the runs in runs/train/exp directory

6. test the model on the test images to label the rest
python detect.py --weights runs/train/exp/weights/best.pt --img 640 --conf 0.25 --source data/images/test


