import pandas as pd
import geopandas as gpd

def main(geojson_path, csv_path):
    """
    Main function to save the geojson file to a csv file
    """
    gdf = gpd.read_file(geojson_path)

    gdf['longitude'] = gdf.geometry.x
    gdf['latitude'] = gdf.geometry.y

    df = pd.DataFrame(gdf.drop(columns=['geometry']))

    df.to_csv(csv_path, index=False)
    print(f"Data saved to {csv_path}")

if __name__ == "__main__":
    main('disabled-parking-bays-gen-2021-wgs84.geojson', 'accessible_parking_dublin.csv')