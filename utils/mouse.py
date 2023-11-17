import pyautogui
from time import sleep


def mouse(data):
    screen_width, screen_height = pyautogui.size()
    mouse_array = data['mousearray']
    for point in mouse_array:
        x = round(point['x']*screen_width, 0)
        y = round(point['y']*screen_height, 0)
        print(f"x {round(point['x']*screen_width, 0)}, y {round(point['y']*screen_height, 0)}")
        pyautogui.moveTo(round(point['x']*screen_width, 0), round(point['y']*screen_height, 0), _pause=False)
        print(round(point['time'], 2)/1000)
        sleep(round(point['time'], 2)/1000)