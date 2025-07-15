#include "serverEvents.h"
#include "publish.h"
#include <Arduino.h>
#include <pinconfig.h>
#include <globals.h>
#include <audio/audio.h>

void test_event_handler(NATS::msg msg)
{
    Serial.printf("Received test event: %s\n", msg.data);
}

void play_audio_handler(NATS::msg msg) {
    String audioPath = String(msg.data);
    Serial.printf("Playing audio from: %s\n", audioPath.c_str());
    playAudioAsync(audioPath);
}

void take_image_handler(NATS::msg msg) {
    Serial.println("Taking image...");
    send_image();    
}

// MESSAGE FORMAT:
// orange on
// orange off
// red on
// red off
void led_handler(NATS::msg msg) {
    String command = String(msg.data);
    if (command == "orange on") {
        digitalWrite(YELLOW_LED_PIN, HIGH);
    } else if (command == "orange off") {
        digitalWrite(YELLOW_LED_PIN, LOW);
    } else if (command == "red on") {
        digitalWrite(RED_LED_PIN, HIGH);
    } else if (command == "red off") {
        digitalWrite(RED_LED_PIN, LOW);
    } else {
        Serial.printf("Unknown LED command: %s\n", command.c_str());
    }
}

void subscribe_to_events()
{
    nats.subscribe("test_event", test_event_handler, NULL, 0);
    nats.subscribe("esp32.play_audio", play_audio_handler, NULL, 0);
    nats.subscribe("esp32.led", led_handler, NULL, 0);
    nats.subscribe("esp32.take_image", take_image_handler, NULL, 0);
}