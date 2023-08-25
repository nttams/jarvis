FROM ubuntu:20.04
WORKDIR /app
CMD ./jarvis

# docker run -d -p 8080:8080 \
#     --restart always \
#     --name=jarvis \
#     --volume "./:/app" \
#     jarvis:latest