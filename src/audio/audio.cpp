#include "audio.h"
#include <Audio.h>
#include <pinconfig.h>
#include <globals.h>

Audio *audio = nullptr;
unsigned long startTimeTimeout = 0;
const unsigned long audioTimeout = 30000; // 30 seconds timeout for audio playback

void playAudio(const String& path) {
    if (audio != nullptr && audio->isRunning()) {
        Serial.println("Audio is already playing.");
        return;
    }

    String url = "http://" + String(SERVER_HOST) + ":" + String(HTTP_PORT) + path;

    audio = new Audio();
    audio->setPinout(I2S_BCLK, I2S_LRC, I2S_DOUT);
    audio->setVolume(AUDIO_VOLUME);
    audio->connecttohost(url.c_str());

    startTimeTimeout = millis();
    while (1) {
        audio->loop();
        if (millis() - startTimeTimeout > audioTimeout) {
            Serial.println("Audio playback timeout reached.");
            break;
        }
        if(!audio->isRunning()) {
            Serial.println("Audio playback finished.");
            break;
        }
    }
    audio->stopSong();
    delete audio;
}

void playAudioAsync(const String& path) {
    if (audio != nullptr && audio->isRunning()) {
        Serial.println("Audio is already playing.");
        return;
    }

    xTaskCreate([](void *arg) {
        String *audioPath = static_cast<String *>(arg);
        playAudio(*audioPath);
        delete audioPath;
        vTaskDelete(nullptr);
    }, "AudioTask", 8192, new String(path), 1, nullptr);
}

bool isAudioPlaying() {
    if (audio != nullptr) {
        return audio->isRunning();
    }
    return false;
}