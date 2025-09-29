#include <Arduino.h>
#include <BluetoothSerial.h>
#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>

// Initialize Bluetooth and MPU6050 objects
BluetoothSerial SerialBT;
Adafruit_MPU6050 mpu;

void setup() {
    // Initialize serial communication for debugging
    Serial.begin(115200);
    Serial.println("Starting Bluetooth Classic setup...");

    // Initialize Bluetooth Classic
    if (!SerialBT.begin("ESP32_DrumLite_Blue")) {  // Set Bluetooth device name
        Serial.println("An error occurred initializing Bluetooth.");
        while (1);
    }
    Serial.println("Bluetooth Classic initialized. Waiting for connection...");

    // Initialize MPU6050 sensor
    if (!mpu.begin()) {
        Serial.println("Failed to initialize MPU6050. Check connections!");
        while (1);
    }
    Serial.println("MPU6050 initialized.");

    // Configure MPU6050 settings
    mpu.setAccelerometerRange(MPU6050_RANGE_8_G);
    mpu.setGyroRange(MPU6050_RANGE_500_DEG);
    mpu.setFilterBandwidth(MPU6050_BAND_5_HZ);
}

void loop() {
    // Check if a Bluetooth client is connected
    if (SerialBT.hasClient()) {
        // Retrieve accelerometer data
        sensors_event_t accel, gyro, temp;
        mpu.getEvent(&accel, &gyro, &temp);

        // Pack accelerometer data into an array
        float accelData[3] = {accel.acceleration.x, accel.acceleration.y, accel.acceleration.z};

        // Send data over Bluetooth Classic as raw bytes
        SerialBT.write((uint8_t*)accelData, sizeof(accelData));

        // Debugging output to Serial Monitor
        Serial.printf("Data sent: X=%.2f, Y=%.2f, Z=%.2f\n", accelData[0], accelData[1], accelData[2]);
    }

    // Delay to match the desired data rate (100 Hz)
    delay(10);
}