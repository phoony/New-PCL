#include "camera.h"
#include <Arduino.h>
#include <Base64.h>

uint8_t imageBuffer[MAX_IMAGE_SIZE];
size_t imageIndex = 0;
String imageBase64;

bool imageBufferEndsWithEndMarker()
{
    if (imageIndex < END_MARKER_LEN)
        return false;
    for (size_t i = 0; i < END_MARKER_LEN; ++i)
    {
        if (imageBuffer[imageIndex - END_MARKER_LEN + i] != END_MARKER[i])
        {
            return false;
        }
    }
    return true;
}

uint8_t *takeImage()
{
    Serial2.println("capture");
    memset(imageBuffer, 0, MAX_IMAGE_SIZE);
    imageIndex = 0;
    unsigned long start = millis();
    unsigned long timeDifference = 0;
    while (!imageBufferEndsWithEndMarker()) {
        timeDifference = millis() - start;
        if(timeDifference > imageTimeout) {
            Serial.println("Image Capture Error: Timeout reached");
            return nullptr;
        }
        if(imageIndex >= MAX_IMAGE_SIZE) {
            Serial.println("Image Capture Error: Buffer overflow");
            return nullptr;
        }

        if (Serial2.available() > 0) {
            imageBuffer[imageIndex++] = Serial2.read();
        }
    }
    return imageBuffer;
}

String *takeImageAsBase64()
{
    uint8_t *image = takeImage();
    imageBase64.clear();
    imageBase64.reserve(MAX_IMAGE_SIZE * 2); // Base64 encoding can increase size

    if (image == nullptr)
    {
        Serial.println("Failed to capture image");
        return nullptr;
    }

    size_t imageSize = imageIndex - END_MARKER_LEN; // Exclude the end marker
    if (imageSize == 0)
    {
        Serial.println("Captured image is empty");
        return nullptr;
    }

    imageBase64 = base64::encode(image, imageSize);
    return &imageBase64;
}