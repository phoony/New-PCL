#include <Arduino.h>
#include <ArduinoJson.h>
#include <pinconfig.h>
#include <audio/audio.h>
#include "globals.h"
#include "camera/camera.h"
#include "publish.h"

SensorData sensorData = {0, false, false};
SensorData lastSensorData = {0, false, false};

Metadata new_metadata(bool isPeriodic)
{
    Metadata metadata;
    metadata.timestamp = time(nullptr);
    metadata.isPeriodic = isPeriodic;
    return metadata;
}

void publish(const char *topic, const char *data)
{
    if (nats.connected)
    {
        nats.publish(topic, data);
    }
    else
    {
        Serial.println("NATS not connected, cannot publish data");
    }
}

void publish_sensor_data(Metadata metadata)
{
    JsonDocument doc;
    doc["timestamp"] = metadata.timestamp;
    doc["isPeriodic"] = metadata.isPeriodic;
    doc["lightSensor"] = sensorData.lightSensor;
    doc["button_1_pressed"] = sensorData.button_1_pressed;
    doc["button_2_pressed"] = sensorData.button_2_pressed;
    doc["isPlayingAudio"] = sensorData.isPlayingAudio;

    char buffer[256];
    serializeJson(doc, buffer);
    publish(SENSOR_TOPIC, buffer);
}

void fetch_sensor_data()
{
    lastSensorData = sensorData; // Store last sensor data for comparison

    sensorData.lightSensor = 0;
    sensorData.button_1_pressed = digitalRead(BUTTON_1_PIN) == LOW;
    sensorData.button_2_pressed = digitalRead(BUTTON_2_PIN) == LOW;
    sensorData.isPlayingAudio = isAudioPlaying();
}

bool sensors_have_changed()
{
    const int32_t threshold = 300; // threshold for light sensor change
    if (abs(sensorData.lightSensor - lastSensorData.lightSensor) >= threshold ||
        sensorData.button_1_pressed != lastSensorData.button_1_pressed ||
        sensorData.button_2_pressed != lastSensorData.button_2_pressed ||
        sensorData.isPlayingAudio != lastSensorData.isPlayingAudio)
    {
        return true;
    }

    return false;
}

void send_image()
{
    String* base64Image = takeImageAsBase64();

    if (base64Image == nullptr)
    {
        Serial.println("Failed to take image");
        return;
    }

    publish("esp32.image", base64Image->c_str());
}

void handle_input() {
    if (sensorData.button_1_pressed && !lastSensorData.button_1_pressed) {
        send_image();
    }
}
