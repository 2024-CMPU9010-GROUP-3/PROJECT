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
      
def send_coach_parking(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['Longitude'], row['Latitude']]
            },
            "type": "coach_parking",
            "details": {
                "number_spaces_avaliable": row['Spaces available']
            }
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
      
def send_dublinbikes(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['Longitude'], row['Latitude']]
            },
            "type": "bike_sharing_station",
            "details": {
                "number_bikes_avaliable": row['Number']
            }
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
      
def send_bike_stands(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['X'], row['Y']]
            },
            "type": "bike_stand",
            "details": {
                "number_stands": row['no_stands']
            }
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
      
def send_water_fountain(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['Long'], row['Lat']]
            },
            "type": "drinking_water_fountain",
            "details": {
            }
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
      
def send_public_toilets(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['Long'], row['Lat']]
            },
            "type": "public_toilet",
            "details": {
            }
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
      
def send_accessible_parking_dlr(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['X'], row['Y']]
            },
            "type": "accessible_parking",
            "details": {
            }
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
      
def send_bike_stands_dlr(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['long'], row['lat']]
            },
            "type": "bike_stand",
            "details": {
                "number_stands": row['nostands']
            }
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
      
def send_bike_stands_south_dublin(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['X'], row['Y']]
            },
            "type": "bike_stand",
            "details": {
            }
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
      
def send_accessible_parking_south_dublin(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['Longitude'], row['Latitude']]
            },
            "type": "accessible_parking",
            "details": {
                "number_spaces": row['No_Of_Spaces']
            }
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
      
def send_public_wifi(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['LONGITUDE'], row['LATITUDE']]
            },
            "type": "public_wifi_access_point",
            "details": {
            }
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
      
def send_libraries(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['Longitude'], row['Latitude']]
            },
            "type": "library",
            "details": {
            }
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
      
def send_accessible_parking(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['longitude'], row['latitude']]
            },
            "type": "accessible_parking",
            "details": {
            }
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
      
def send_multistorey_car_park(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['Longitude'], row['Latitude']]
            },
            "type": "multistorey_car_parking",
            "details": {
            }
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
      
def send_bleeperbike(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['longitude'], row['latitude']]
            },
            "type": "bike_sharing_station",
            "details": {
            }
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
      
def send_parking_meter(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['longitude'], row['latitude']]
            },
            "type": "parking_meter",
            "details": {
            }
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
      
def send_public_bins(session, url, row):
    point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['longitude'], row['latitude']]
            },
            "type": "public_bins",
            "details": {
            }
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
    send_parking_spots_to_api(data)
    

if __name__ == "__main__":
    main(sys.argv[1])
