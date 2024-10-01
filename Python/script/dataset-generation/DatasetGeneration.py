import requests
from PIL import Image
from io import BytesIO
import os
import csv
import pandas as pd


df = pd.read_csv(
    "dataset-generation/sport-pitches-and-facilities.csv")
location_data = df.to_dict('records')


def get_images(name, latitude, longitude):
    # Define the API endpoint URL
    url = f'https://api.mapbox.com/styles/v1/mapbox/satellite-v9/static/{longitude},{latitude},18,0,0/400x400?access_token=pk.eyJ1Ijoia2F1c3R1Ymh0cml2ZWRpIiwiYSI6ImNtMWo2NndsbzB4N3EycHM1aGF2cDd5NzkifQ.4aegzX6Kfy3zW8pHkLWU7Q'
    response = requests.get(url)
    print(response)
    if response.status_code == 200:
        # Step 3: Open the image from the response content
        img = Image.open(BytesIO(response.content))
        if not os.path.exists('output'):
            os.makedirs('output')
        # Step 4: Define the output folder path
        output_folder = 'output'

        # Step 6: Define the path where the image will be saved
        output_path = os.path.join(output_folder, f'{name}.png')

        # Step 7: Save the image to the output folder
        img.save(output_path)

        print(f"Image saved to {output_path}")


for place in location_data:
    name = place["Title"]
    if ("/" in name):
        name = name.replace("/", "-")
    latitude = place["Latitude"]
    longitude = place["Longitude"]
    get_images(name, latitude, longitude)
