### Instructions
I did this in VS code, the code is to be run in the terminal 

1. Add all images and labels to the "images" & "labels" folders

2. Run the "data split" file to split the image data into train & validation
python data_split.py

3. Clone the Yolov5 GitHub repo
git clone https://github.com/ultralytics/yolov5.git  

4. Cd into the cloned directory & fulfill requirements to run the project
cd yolov5
pip install -r requirements.txt

5.1 Add the data.yaml and hyp.yaml files to the yolov5 directory

6. Download pre-trained weights for the yolov5 model (i've chosen the small model - yolov5s.pt -- to start with)
python detect.py --weights yolov5s.pt --source data/images

7. fine tune model by training & validating on our images
you can play around with the parameters in hyp.yaml and the number of epochs
python train.py --img 640 --batch 16 --epochs 10 --data data.yaml --weights C:\{YOUR FILE PATH}\yolov5\runs\train\exp\weights\best.pt --hyp hyp.yaml --cache

8. view results of the runs in runs/train/exp directory

9. test the model on the test images to label the rest
python detect.py --weights runs/train/exp8/weights/best.pt --img 640 --conf 0.25 --source C:\{YOUR FILE PATH}\dataset_ann\images\test


