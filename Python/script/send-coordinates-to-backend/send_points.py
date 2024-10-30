import pandas as pd
import requests

def send_parking_spots_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
        point_data = {
            "longlat": {
                "type": "Point",
                "coordinates": [row['longitude'], row['latitude']]
            },
            "type": "parking",
            "details": {
            }
        }
        
        try:
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['longitude']} and latitude {row['latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_coach_parking_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['Longitude']} and latitude {row['Latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_dublinbikes_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['Longitude']} and latitude {row['Latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_bike_stands_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['X']} and latitude {row['Y']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_water_fountain_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['Long']} and latitude {row['Lat']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_public_toilets_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['Long']} and latitude {row['Lat']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_accessible_parking_dlr_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['X']} and latitude {row['Y']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_bike_stands_dlr_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['long']} and latitude {row['lat']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_bike_stands_south_dublin_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['X']} and latitude {row['Y']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_accessible_parking_south_dublin_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['Longitude']} and latitude {row['Latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_public_wifi_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['LONGITUDE']} and latitude {row['LATITUDE']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")


def send_libraries_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['Longitude']} and latitude {row['Latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_accessible_parking_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['longitude']} and latitude {row['latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_multistorey_car_park_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['Longitude']} and latitude {row['Latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_bleeperbike_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['longitude']} and latitude {row['latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_parking_meter_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['longitude']} and latitude {row['latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")

def send_public_bins_to_api(data):
    """
    Function to send the data as post requests to the backend
    """
    url = "http://localhost:8080/v1/private/points/"

    successful_posts = 0

    for index, row in data.iterrows():
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
            response = requests.post(url, json=point_data)
            
            if response.status_code == 200 or response.status_code == 201:
                print(f"Successfully added point with longitude {row['longitude']} and latitude {row['latitude']}")
                successful_posts += 1
            else:
                print(f"Failed to add point: {response.status_code}, {response.text}")
        
        except Exception as e:
            print(f"Error sending data to API: {e}")
    
    print(f"Total successful posts: {successful_posts}/{len(data)}")


def read_csv_file(csv_file_path):
    """
    Functiion to read the csv file in a pandas dataframe
    """
    try:
        data = pd.read_csv(csv_file_path)
        return data
    except Exception as e:
        print(f"Error reading CSV file: {e}")
        return None

def main(csv_file_path):
    """
    Main function to send all the coordiantes in a csv file to the backend
    """
    data = read_csv_file(csv_file_path)

    if data is not None:
        send_coach_parking_to_api(data)
        #send_dublinbikes_to_api(data)
        #send_bike_stands_to_api(data)
        #send_water_fountain_to_api(data)
        #send_public_toilets_to_api(data)
        #send_accessible_parking_dlr_to_api(data)
        #send_bike_stands_dlr_to_api(data)
        #send_bike_stands_south_dublin_to_api(data)
        #send_accessible_parking_south_dublin_to_api(data)
        #send_public_wifi_to_api(data)
        #send_libraries_to_api(data)
        #send_accessible_parking_to_api(data)
        #send_multistorey_car_park_to_api(data)
        #send_bleeperbike_to_api(data)
        #send_parking_meter_to_api(data)
        #send_public_bins_to_api(data)


if __name__ == "__main__":
    main('2020-coach-parking-dcc-1.csv')
    #main('dublinbikes.csv')
    #main('bike_stands.csv')
    #main('drinking_water_fountain.csv')
    #main('public-toilets-dcc-2021-updated.csv')
    #main('accessible_parking_bays_dlr2016.csv')
    #main('bicycleparkingstandsdlr.csv')
    #main('Bicycle_Parking_Stands_SDCC.csv')
    #main('Accessible_Parking_Spaces_SDCC.csv')
    #main('dcc_wifi4eu_locations.csv')
    #main('dublincitylibrarylocation2021.csv')
    #main('accessible_parking_dublin.csv')
    #main('multi_story_car_parks_location.26052021.csv')
    #main('bleeperbike_map.csv')
    #main('parking-meter.csv')
    #main('public_bin.csv)

