#ifndef GLOBALS_H
#define GLOBALS_H

#include <WiFi.h>
#include <ArduinoNATS.h>

#define WIFI_SSID ""
#define WIFI_PASSWORD ""
#define SERVER_HOST ""

#define NATS_PORT 4222
#define HTTP_PORT 5000

#define AUDIO_VOLUME 10

extern WiFiClient net;
extern NATS nats;

#endif // GLOBALS_H