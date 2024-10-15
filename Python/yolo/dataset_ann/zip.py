import os
import zipfile
import sys

def zip_directories():
    current_dir = os.getcwd()
    for item in os.listdir(current_dir):
        item_path = os.path.join(current_dir, item)
        if os.path.isdir(item_path):
            zip_file = f"{item}.zip"
            with zipfile.ZipFile(zip_file, 'w', zipfile.ZIP_DEFLATED) as zipf:
                for root, dirs, files in os.walk(item_path):
                    for file in files:
                        file_path = os.path.join(root, file)
                        arcname = os.path.relpath(file_path, start=item_path)
                        zipf.write(file_path, arcname)

def unzip_files():
    current_dir = os.getcwd()
    for item in os.listdir(current_dir):
        if item.endswith('.zip'):
            zip_file_path = os.path.join(current_dir, item)
            with zipfile.ZipFile(zip_file_path, 'r') as zipf:
                zipf.extractall(os.path.splitext(zip_file_path)[0])

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python zip.py <zip|unzip>")
        sys.exit(1)
    
    operation = sys.argv[1].lower()
    if operation == 'zip':
        zip_directories()
    elif operation == 'unzip':
        unzip_files()
    else:
        print("Invalid argument. Use 'zip' to zip directories or 'unzip' to unzip files.")
        sys.exit(1)