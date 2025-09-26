import threading
import serial
import struct
import time

# Replace MAC addresses with COM ports for Bluetooth Classic
BLUE_COM_PORT = "COM13"  # Replace with the correct COM port for the blue drumstick
GREEN_COM_PORT = "COM11"  # Replace with the correct COM port for the green drumstick
BAUD_RATE = 115200  # Match with ESP32 firmware settings

# Global variables for accelerometer data
average_accel_blue = 0
average_accel_green = 0
sign_blue = 0
sign_green = 0


def read_accel_data(serial_conn, color):
    """Read and process accelerometer data from the drumstick."""
    global average_accel_blue, average_accel_green, sign_blue, sign_green

    try:
        if serial_conn.in_waiting >= 12:  # Each float is 4 bytes, so expect 12 bytes for 3 floats
            data = serial_conn.read(12)  # Read 12 bytes
            accel_x, accel_y, accel_z = struct.unpack('<fff', data)

            if color == "blue":
                sign_blue = 1 if accel_z >= 0 else -1
                average_accel_blue = sign_blue * (abs(accel_x) + abs(accel_z)) / 2
            elif color == "green":
                sign_green = 1 if accel_z >= 0 else -1
                average_accel_green = sign_green * (abs(accel_x) + abs(accel_z)) / 2

    except Exception as e:
        print(f"Error reading data from {color} drumstick: {e}")


def connect_to_drumstick(color, com_port):
    """Handle serial connection for a specific drumstick."""
    try:
        conn = serial.Serial(com_port, BAUD_RATE, timeout=1)
        print(f"{color.capitalize()} drumstick connected on {com_port}")
        while True:
            read_accel_data(conn, color)
            time.sleep(0.01)  # Match drumstick data rate
    except serial.SerialException as e:
        print(f"Error connecting to {color} drumstick on {com_port}: {e}")


def get_accel_blue():
    """Return the latest processed value for the blue drumstick."""
    global average_accel_blue
    return average_accel_blue


def get_accel_green():
    """Return the latest processed value for the green drumstick."""
    global average_accel_green
    return average_accel_green


def start_bt_threads():
    """Start threads for both drumsticks."""
    blue_thread = threading.Thread(target=connect_to_drumstick, args=("blue", BLUE_COM_PORT))
    green_thread = threading.Thread(target=connect_to_drumstick, args=("green", GREEN_COM_PORT))
    blue_thread.start()
    green_thread.start()


if __name__ == "__main__":
    start_bt_threads()

    try:
        while True:
            print(f"Blue: {get_accel_blue():.2f}, Green: {get_accel_green():.2f}")
            time.sleep(0.5)  # Adjust display rate as needed
    except KeyboardInterrupt:
        print("Shutting down.")