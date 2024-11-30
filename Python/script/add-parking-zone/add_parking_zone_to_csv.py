import geopandas as gpd
import pandas as pd
import matplotlib.pyplot as plt

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

    zones_gdf = zones_gdf.set_crs(epsg=4326)
    parking_gdf = parking_gdf.set_crs(epsg=4326)


    print(parking_gdf.crs)
    print(zones_gdf.crs)
    print(zones_gdf.geometry.is_valid)


    if 'index_right' in parking_gdf.columns:
        parking_gdf = parking_gdf.drop(columns='index_right')
    if 'index_right' in zones_gdf.columns:
        zones_gdf = zones_gdf.drop(columns='index_right')

    parking_with_zones = gpd.sjoin(parking_gdf, zones_gdf, how="left", predicate="within")

    print(parking_with_zones.columns)


    fig, ax = plt.subplots(figsize=(10, 10))
    zones_gdf.plot(ax=ax, color='lightblue', edgecolor='black', alpha=0.5)
    parking_gdf.plot(ax=ax, color='red', markersize=5)
    plt.show()

    zones_gdf.plot(facecolor='lightblue', edgecolor='black', alpha=0.5)
    plt.show()

    base = zones_gdf.plot(facecolor='lightblue', edgecolor='black', alpha=0.5)
    parking_gdf.plot(ax=base, color='red', markersize=10)
    plt.show()



    print("Parking Spots Bounding Box:")
    print(parking_gdf.total_bounds)  # MinX, MinY, MaxX, MaxY for parking spots

    print("Zones Bounding Box:")
    print(zones_gdf.total_bounds) 

    parking_with_zones['parking_zone'] = parking_with_zones['zone']
    parking_with_zones['parking_zone_tarif'] = parking_with_zones['tarif']

    parking_with_zones.drop(columns='geometry', inplace=True)
    parking_with_zones.to_csv(parking_csv, index=False)

if __name__ == "__main__":
    main("coordinates_in_-6.3169_53.3767--6.3159_53.3773.csv", "tarif_zones.geojson")
    #main("coordinates_in_-6.2708_53.3456--6.27_53.3459.csv", "tarif_zones.geojson")
