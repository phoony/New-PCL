#include "globals.h"

WiFiClient net;
NATS nats(&net, SERVER_HOST, NATS_PORT);