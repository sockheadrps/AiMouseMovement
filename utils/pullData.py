from pymongo import MongoClient
from dotenv import load_dotenv
import os
import json

file_path = 'output.json'
dotenv_path = os.path.join(os.getcwd(), 'server', '.env')
load_dotenv(dotenv_path)

MONGO_URI = os.getenv("MONGO_URI")

client = MongoClient(MONGO_URI)

# # Access a specific collection
collection = client.mousedb.mouse
data_points = list(collection.find())

def get_file_size(file_path):
    size_in_bytes = os.path.getsize(file_path)
    size_in_kb = size_in_bytes / 1024
    size_in_mb = size_in_kb / 1024

    return size_in_mb

if __name__ == "__main__":
    with open(file_path, 'w') as file:
        json.dump(data_points, file, default=str)

    print(f'File Size: {get_file_size(file_path):.2f} MB)')