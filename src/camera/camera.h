#ifndef CAMERA_H
#define CAMERA_H

#include <Arduino.h>

// --- Camera setup ---
const unsigned long imageTimeout = 5000;
const uint8_t END_MARKER[] = { 0xFF, 0xD9, 0x00, 0x00, 0xDE, 0xAD, 0xBE, 0xEF };
const size_t END_MARKER_LEN = sizeof(END_MARKER);
const int MAX_IMAGE_SIZE = 100 * 256;

uint8_t* takeImage();
String* takeImageAsBase64();

#endif // CAMERA_H