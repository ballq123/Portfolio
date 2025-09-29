# DrumLite – Bluetooth Telemetry & Drumstick Detection

**Timeline:** Aug 2024 – Dec 2024  
**Team Size:** 3 (Capstone Project, CMU ECE)
**Focus:** Embedded firmware + wireless communication + computer vision  
**Tech Stack:** ESP32 (C/Arduino), Bluetooth Classic, Adafruit MPU-6050, Python (OpenCV, threading, serial)

---

## Overview
DrumLite is a **portable electronic drum kit** designed for accessibility and portability.  
It integrates **ESP32 microcontrollers with accelerometer sensors** and a **computer vision module** for drumstick detection. Drum hits are detected in real time and transmitted via Bluetooth to trigger audio playback with an **average end-to-end latency of ~52ms**.  

This project demonstrates skill in **embedded systems, wireless communication, and computer vision**.
This submodule demonstrates two key components of DrumLite:
	1.	ESP32 Bluetooth Telemetry – Streams accelerometer readings from an ESP32 over Bluetooth Classic (SPP) at ~100 Hz.
	2.	Drumstick Tip Detection – A color-segmentation OpenCV pipeline that finds blue/green drumstick tips and maps them to drum pads.
Together, these enable real-time hit detection with ~52 ms end-to-end latency (Bluetooth ~17 ms, CV ~30 ms, accelerometer ~5 ms).


---

## Files

- [`main.cpp`](code-snippets/bluetooth/src/main.cpp) → ESP32 firmware  
  - Initializes Bluetooth Classic and streams raw accelerometer data (`x, y, z` floats) at ~100 Hz.  
  - Demonstrates **embedded programming**, **sensor integration**, and **low-latency wireless data transmission**.  

- [`client.py`](code-snippets/bluetooth/client.py) → Host-side Bluetooth client  
  - Opens serial ports to receive ESP32 telemetry, unpacks float data, and computes signed average accelerations.  
  - Demonstrates **multithreading**, **binary protocol handling (struct unpack)**, and **real-time data processing** in Python.  

- [`drumstick_tracking.py`](code-snippets/drumstick_tracking.py) → Computer vision module  
  - Uses OpenCV color segmentation (blue/green HSV thresholds) to detect drumstick tips and map them to four drum pads.  
  - Demonstrates **image preprocessing (CLAHE, Gaussian blur)**, **contour detection**, and **geometric mapping** from vision to physical drum pads.  

---

## Impact
- Achieved **~52 ms end-to-end latency** from drumstick hit to audio playback, exceeding our original 100 ms goal by over 2x and reaching responsiveness comparable to commercial e-drums  
- Enabled **reliable audio playback across four drum faces**, validated through user testing and threshold tuning that reduced false positives  
- Selected as the **only ECE capstone project** to showcase at Carnegie Mellon’s TechSpark Engineering Exposition alongside teams from MechE, BME, and other engineering departments  
- Awarded **“Most Engaging Project”** at the TechSpark Expo (Fall 2024), highlighting both technical excellence and user accessibility focus  
- Presented at the **CMU ECE Senior Design Expo**, demonstrating accessibility-driven innovation to faculty, peers, and industry reviewers

---

## Skills Demonstrated
- **Embedded Systems:**  
  - ESP32 firmware development in C/Arduino  
  - Sensor interfacing with MPU-6050 accelerometer  
  - Real-time data streaming at fixed sampling rates  

- **Wireless Systems:**  
  - Bluetooth Classic (SPP) integration on ESP32  
  - Host–device synchronization at ~100 Hz  
  - Efficient binary data packing/unpacking  

- **Systems & Concurrency:**  
  - Multithreading in Python for parallel serial connections  
  - Real-time buffering and signal processing  
  - Fault handling in long-lived data streams  

- **Computer Vision:**  
  - OpenCV color segmentation with adjustable HSV thresholds  
  - Contrast enhancement (CLAHE) for robust detection under variable lighting  
  - Contour-based detection and spatial mapping to defined drum pads 
