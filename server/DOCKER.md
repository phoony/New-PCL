# Smart Fridge
## Building
```bash
docker build -t smart-fridge .
```
## Running
```bash
docker run --name smart-fridge-server  --restart unless-stopped -p 5000:5000 -p 4222:4222 -e HTTP_PORT=5000 -e NATS_PORT=4222 -e OPENAI_API_KEY=bTVvV_qaX8g60lF8KMV7UDa8XhxPS_OO0u9sGq6dWjo4Q7qLrncgVFj143NzYSu8sYwoXYJQlDT3BlbkFJJ8Fr_SPlZjUWKrRiGRRq76oH6du2Couw3CaoAEmzyk59EEoktwZIJY3NKBq0lZfGoU_392LbgA smart-fridge
```

## Environment Variables
- `HTTP_PORT`
- `NATS_PORT`
- `OPENAI_API_KEY`
