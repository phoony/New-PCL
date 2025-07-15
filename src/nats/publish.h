#ifndef PUBLISH_H
#define PUBLISH_H

#include <time.h>

#define SENSOR_TOPIC "esp32.sensors"

struct SensorData
{
    int32_t lightSensor;
    bool button_1_pressed;
    bool button_2_pressed;
    bool isPlayingAudio;
};

struct Metadata
{
    time_t timestamp;
    bool isPeriodic;
    bool isPlayingAudio;
};

struct SensorMessage
{
    SensorData sensorData;
    Metadata metadata;
};

Metadata new_metadata(bool isPeriodic);

void publish_sensor_data(const char *topic, const char *data);
void publish_sensor_data(Metadata metadata);

void fetch_sensor_data();
bool sensors_have_changed();
void handle_input();
void send_image();

#endif // PUBLISH_H