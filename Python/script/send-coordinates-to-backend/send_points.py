import pandas as pd
import requests
import argparse
import numpy as np
import os
import re
import glob
from tqdm import tqdm
from concurrent.futures import ThreadPoolExecutor, as_completed


def infer_type_from_filename(filename, type_list):
    """
    Infer the point types from the filename based on the Hamming distance to the types in type_list.
    Returns a list of types sorted by increasing Hamming distance.
    """
    basename = os.path.basename(filename).lower()
    # Remove file extension
    name_without_ext = os.path.splitext(basename)[0]
    # Remove non-alphanumeric characters, spaces, numbers, and underscores
    name_cleaned = re.sub(r'[\W\d_]+', '', name_without_ext)
    # Now compute Hamming distance to each type in the list
    distances = []
    for point_type in type_list:
        # Clean the point_type similarly
        type_cleaned = re.sub(r'[\W\d_]+', '', point_type.lower())
        # Pad the shorter string with spaces to make them equal length
        max_len = max(len(name_cleaned), len(type_cleaned))
        s1 = name_cleaned.ljust(max_len)
        s2 = type_cleaned.ljust(max_len)
        # Compute Hamming distance
        distance = sum(ch1 != ch2 for ch1, ch2 in zip(s1, s2))
        distances.append((distance, point_type))
    # Sort the list by distance
    distances.sort()
    # Return the sorted list of types
    return distances


def send_parking_spot(session, url, row, point_type):
    point_data = {
        "longlat": {
            "type": "Point",
            "coordinates": [row['longitude'], row['latitude']]
        },
        "type": point_type,
        "details": {}
    }
    # Build details excluding 'latitude' and 'longitude'
    details = row.drop(['latitude', 'longitude']).replace(
        {np.nan: None}).to_dict()
    point_data['details'] = details
    try:
        response = session.post(url, json=point_data)
        if response.status_code in (200, 201):
            return True
        else:
            print(
                f"Failed to add point: {response.status_code}, {response.text}")
            return False
    except Exception as e:
        print(f"Error sending data to API: {e}")
        return False


def send_parking_spots_to_api(data, point_type, url):
    """
    Function to send the data as POST requests to the backend using multithreading.
    """
    successful_posts = 0
    total = len(data)

    # Create a shared session
    session = requests.Session()

    with ThreadPoolExecutor() as executor:
        futures = []
        for index, row in data.iterrows():
            futures.append(executor.submit(
                send_parking_spot, session, url, row, point_type))

        with tqdm(total=total, desc="Uploading points", unit="point") as pbar:
            for future in as_completed(futures):
                result = future.result()
                if result:
                    successful_posts += 1
                pbar.update(1)
                pbar.set_postfix(successful=successful_posts,
                                 failed=(pbar.n - successful_posts))

    session.close()


def read_csv_file(csv_file_path):
    """
    Function to read the CSV file into a pandas DataFrame.
    """
    try:
        data = pd.read_csv(csv_file_path)
        # Standardize column names
        data.columns = [col.strip().lower() for col in data.columns]
        # Convert latitude and longitude to numeric
        data['latitude'] = pd.to_numeric(data['latitude'], errors='coerce')
        data['longitude'] = pd.to_numeric(data['longitude'], errors='coerce')
        # Drop rows with NaN latitude or longitude
        data = data.dropna(subset=['latitude', 'longitude'])
        return data
    except Exception as e:
        print(f"Error reading CSV file: {e}")
        return None


def main():
    """
    Main function to send all the coordinates in CSV files to the backend.
    """
    parser = argparse.ArgumentParser(
        description='Send points to API.',
        formatter_class=argparse.RawTextHelpFormatter,
        epilog="""Examples:
      python script_name.py data.csv --type parking --env prod
          - Sends points from 'data.csv' with type 'parking' to the production environment.

      python script_name.py all --type auto --env dev
          - Processes all CSV files in the current directory.
          - Infers point types from filenames.
          - Sends data to the development environment.
          \n
    """
    )
    parser.add_argument(
        'csv_file_path',
        help='Path to the CSV file or "all" to process all CSV files in the directory.'
    )
    parser.add_argument(
        '--type',
        required=True,
        help='Type of the points. Use "auto" to infer from filename.'
    )
    parser.add_argument(
        '--env',
        choices=['prod', 'dev'],
        default='dev',
        help='Environment to use: "prod" for production, "dev" for development (default: dev).'
    )
    args = parser.parse_args()

    # Determine the API URL based on the environment
    if args.env == 'prod':
        api_url = "http://localhost:8080/v1/private/points/"
    else:
        api_url = "http://localhost:8081/v1/private/points/"

    # List of possible types
    type_list = [
        'parking',
        'unknown',
        'coach_parking',
        'bike_sharing_station',
        'bike_stand',
        'drinking_water_fountain',
        'public_toilet',
        'accessible_parking',
        'public_wifi_access_point',
        'library',
        'multistorey_car_parking',
        'parking_meter',
        'public_bins'
    ]

    if args.csv_file_path.lower() == 'all':
        # Process all CSV files in the directory
        csv_files = glob.glob('*.csv')
        if not csv_files:
            print("No CSV files found in the directory.")
            return
        for csv_file in csv_files:
            print(f"\nProcessing file: {csv_file}")
            if args.type.lower() == 'auto':
                # Get the sorted list of types based on Hamming distance
                sorted_types = infer_type_from_filename(csv_file, type_list)
                point_type = None  # Initialize point_type
                for distance, proposed_type in sorted_types:
                    # Confirmation dialogue
                    while True:
                        confirmation = input(
                            f"Do you want to proceed with the inferred type '{proposed_type}' for file '{csv_file}'? (y/n): ").strip().lower()
                        if confirmation in ['y', 'yes', '']:
                            # User confirmed the inferred type
                            point_type = proposed_type
                            break
                        elif confirmation in ['n', 'no']:
                            # User did not confirm, try next closest type
                            break
                        else:
                            print(
                                "Invalid input. Please enter 'Y' for yes or 'N' for no.")
                    if point_type:
                        # User confirmed, exit loop
                        break
                if not point_type:
                    # If no type was confirmed, ask user to input type
                    print(
                        f"Could not determine point type for file '{csv_file}'.")
                    print("Available types:")
                    for t in type_list:
                        print(f" - {t}")
                    while True:
                        point_type_input = input(
                            "Please enter the correct type from the list above: ").strip()
                        if point_type_input in type_list:
                            point_type = point_type_input
                            break
                        else:
                            print("Invalid type entered. Please try again.")
            else:
                point_type = args.type

            data = read_csv_file(csv_file)
            if data is not None:
                send_parking_spots_to_api(data, point_type, api_url)
    else:
        # Original code for a single CSV file
        if args.type.lower() == 'auto':
            # Get the sorted list of types based on Hamming distance
            sorted_types = infer_type_from_filename(
                args.csv_file_path, type_list)
            point_type = None  # Initialize point_type
            for distance, proposed_type in sorted_types:
                # Confirmation dialogue
                while True:
                    confirmation = input(
                        f"Do you want to proceed with the inferred type '{proposed_type}'? (y/n): ").strip().lower()
                    if confirmation in ['y', 'yes', '']:
                        # User confirmed the inferred type
                        point_type = proposed_type
                        break
                    elif confirmation in ['n', 'no']:
                        # User did not confirm, try next closest type
                        break
                    else:
                        print(
                            "Invalid input. Please enter 'Y' for yes or 'N' for no.")
                if point_type:
                    # User confirmed, exit loop
                    break
            if not point_type:
                # If no type was confirmed, ask user to input type
                print("Could not determine point type.")
                print("Available types:")
                for t in type_list:
                    print(f" - {t}")
                while True:
                    point_type_input = input(
                        "Please enter the correct type from the list above: ").strip()
                    if point_type_input in type_list:
                        point_type = point_type_input
                        break
                    else:
                        print("Invalid type entered. Please try again.")
        else:
            point_type = args.type

        data = read_csv_file(args.csv_file_path)
        if data is not None:
            send_parking_spots_to_api(data, point_type, api_url)


if __name__ == "__main__":
    main()
