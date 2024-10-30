import pandas as pd
import requests
import sys
from tqdm import tqdm
from concurrent.futures import ThreadPoolExecutor, as_completed

def send_parking_spot(session, url, row):
    point_data = {
        "longlat": {
            "type": "Point",
            "coordinates": [row['longitude'], row['latitude']]
        },
        "type": "parking",
        "details": {}
    }
    try:
        response = session.post(url, json=point_data)
        if response.status_code in (200, 201):
            return True
        else:
            print(f"Failed to add point: {response.status_code}, {response.text}")
            return False
    except Exception as e:
        print(f"Error sending data to API: {e}")
        return False

def send_parking_spots_to_api(data):
    """
    Function to send the data as POST requests to the backend using multithreading.
    """
    url = "http://localhost:8080/v1/private/points/" # This is intentionally hardcoded, this will never change in prod
    successful_posts = 0
    total = len(data)

    # Create a shared session
    session = requests.Session()

    with ThreadPoolExecutor() as executor:
        futures = []
        for index, row in data.iterrows():
            futures.append(executor.submit(send_parking_spot, session, url, row))

        with tqdm(total=total, desc="Uploading parking spots", unit="spot") as pbar:
            for future in as_completed(futures):
                result = future.result()
                if result:
                    successful_posts += 1
                pbar.update(1)
                pbar.set_postfix(successful=successful_posts, failed=(pbar.n - successful_posts))

    session.close()

def read_csv_file(csv_file_path):
    """
    Function to read the CSV file into a pandas DataFrame.
    """
    try:
        data = pd.read_csv(csv_file_path)
        return data
    except Exception as e:
        print(f"Error reading CSV file: {e}")
        return None

def main(csv_file_path):
    """
    Main function to send all the coordinates in a CSV file to the backend.
    """
    data = read_csv_file(csv_file_path)
    if data is not None:
        send_parking_spots_to_api(data)

if __name__ == "__main__":
    main(sys.argv[1])
