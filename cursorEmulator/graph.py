from dotenv import load_dotenv
import os
import pyautogui
from time import sleep
from pymongo import MongoClient
import matplotlib.pyplot as plt


dotenv_path = os.path.join(os.getcwd(), 'server', '.env')
# Load environment variables from .env file
load_dotenv(dotenv_path)

# Access environment variables
MONGO_URI = os.getenv("MONGO_URI")
print(MONGO_URI)

# Connect to the MongoDB server
client = MongoClient(MONGO_URI)


# # Access a specific collection
collection = client.mousedb.mouse
print(collection)
screen_width, screen_height = pyautogui.size()
i = 0
for col in collection.find():
    i+=1
    mouse_array = col['mousearray']
    height = col['windowheight']
    width = col['windowwidth']
    # Extract 'x' and 'y' values from each object
    x_values = [round(point['x']*screen_width, 0) for point in mouse_array]
    y_values = [round(point['y']*screen_height, 0) for point in mouse_array]

    # Create a scatter plot
    plt.scatter(x_values, y_values, label='Data Points')

    # Set labels and title
    plt.xlabel('X')
    plt.ylabel('Y')
    plt.title('Scatter Plot of Data Points')

    # Add legend
    plt.legend()

    # Show the plot
    plt.show()
    if i == 10:
        break






    # for point in mouse_array:
    #     x = round(point['x']*screen_width, 0)
    #     y = round(point['y']*screen_height, 0)
    #     print(f"x {round(point['x']*screen_width, 0)}, y {round(point['y']*screen_height, 0)}")
    #     pyautogui.moveTo(round(point['x']*screen_width, 0), round(point['y']*screen_height, 0), _pause=False)
    #     print(round(point['time'], 2)/1000)
    #     sleep(round(point['time'], 2)/1000)

