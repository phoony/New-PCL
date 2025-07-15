#include <setup/setup.h>
#include <Arduino.h>
#include <WiFi.h>
#include <ArduinoNATS.h>
#include <time.h>
#include "globals.h"
#include "nats/serverEvents.h"

void connect_wifi()
{
	Serial.print("Connecting to WiFi");
	WiFi.mode(WIFI_STA);
	WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
	while (WiFi.status() != WL_CONNECTED)
	{
		Serial.print(".");
		delay(500);
	}

	Serial.println("\nWiFi connected");
}

void setup_time()
{
	Serial.println("Setting up time");

	// Initialize NTP client
	// Set the time zone to UTC+0
	configTime(0, 0, "pool.ntp.org", "time.nist.gov");
	struct tm timeinfo;
	if (!getLocalTime(&timeinfo))
	{
		Serial.println("Failed to obtain time");
		return;
	}
	Serial.println(&timeinfo, "%A, %B %d %Y %H:%M:%S");
}

void nats_connect()
{
	if (nats.connect())
	{
		Serial.println("NATS connected successfully");
		subscribe_to_events();
	}
	else
	{
		Serial.println("Failed to connect to NATS");
	}
}