import geopandas as gpd
import pandas as pd

def main(parking_csv, zones_geojson):
    """
    Assigns the parking zones to the parking spots found 

    Parameters:
    parking_csv: Path to the csv file containing the parking spot coordinates
    zones_geojson: Path to the geojson file containing the parking zones
    """
    parking_df = pd.read_csv(parking_csv)
    parking_gdf = gpd.GeoDataFrame(
        parking_df, geometry=gpd.points_from_xy(parking_df['longitude'], parking_df['latitude'])
    )

    zones_gdf = gpd.read_file(zones_geojson)

    zones_gdf = zones_gdf.to_crs(epsg=4326)
    parking_gdf = parking_gdf.set_crs(epsg=4326)

    if 'index_right' in parking_gdf.columns:
        parking_gdf = parking_gdf.drop(columns='index_right')
    if 'index_right' in zones_gdf.columns:
        zones_gdf = zones_gdf.drop(columns='index_right')

    parking_with_zones = gpd.sjoin(parking_gdf, zones_gdf, how="left", predicate="within")

    print(parking_with_zones.columns)

    parking_with_zones['parking_zone'] = parking_with_zones['zone_left']
    parking_with_zones['parking_zone_tarif'] = parking_with_zones['tarif_left']

    parking_with_zones.drop(columns='geometry', inplace=True)
    parking_with_zones.to_csv(parking_csv, index=False)

if __name__ == "__main__":
    main("coordinates_in_-6.2264_53.4194--6.2219_53.4221.csv", "tarif_zones.geojson")