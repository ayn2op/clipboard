FROM golang:1.24
RUN apt-get update && apt-get install -y \
      xvfb libx11-dev libegl1-mesa-dev libgles2-mesa-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
WORKDIR /app
COPY . .
CMD [ "sh", "-c", "./tests/test-docker.sh" ]
