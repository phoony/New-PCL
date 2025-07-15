#include <WiFi.h>
#include <Audio.h>
#include <HTTPClient.h>
#include <globals.h>
#include <pinconfig.h>
#include <time.h>
#include <nats/serverEvents.h>
#include <nats/publish.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include <setup/setup.h>



void periodicSensors(void *arg)
{
	while (1)
	{
		nats.process();
		//Serial.println("Looping periodic");
		if (nats.connected)
		{
			fetch_sensor_data();
			publish_sensor_data(new_metadata(true));
			//Serial.println("Published periodic sensor data");
		}
		delay(1000);
	}
}

void setup()
{
	Serial.begin(115200);
 // - Light sensor setup -
    pinMode(PIN_LIGHT_SENSOR, OUTPUT);
    digitalWrite(PIN_LIGHT_SENSOR, LOW);

    // - Button setup -
    pinMode(BUTTON_1_PIN, INPUT_PULLUP);
    pinMode(BUTTON_2_PIN, INPUT_PULLUP);

    // - LED setup -
    pinMode(RED_LED_PIN, OUTPUT);
    pinMode(YELLOW_LED_PIN, OUTPUT);

	connect_wifi();
	setup_time();
	nats_connect();
	xTaskCreate(periodicSensors, "PeriodicSensors", 4096, NULL, 2, NULL);
}

static bool flag = false;

void loop()
{
start:

	if (WiFi.status() != WL_CONNECTED)
		connect_wifi();
	if (!nats.connected)
	{
		Serial.println("NATS not connected, trying to reconnect...");
		delay(1000);
		nats_connect();
		goto start;
	}

	fetch_sensor_data();
	if (sensors_have_changed())
	{
		handle_input();
		Metadata metadata = new_metadata(false);
		publish_sensor_data(metadata);
		Serial.println("Published non periodic sensor data");
	}

	nats.process();
	delay(30);
}
