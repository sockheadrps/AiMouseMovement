import matplotlib.pyplot as plt
import pyautogui
import json


data_file = 'output.json'

def read_json_file(file_path):
    with open(file_path, 'r') as file:
        data = json.load(file)
    return data


def graph_data(data):
    screen_width, screen_height = pyautogui.size()
    mouse_array = data['mousearray']
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


json_data = read_json_file(data_file)
graph_data(json_data[0])

