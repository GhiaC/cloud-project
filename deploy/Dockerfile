FROM golang:latest

# go build -ldflags="-w -s" main.go
# docker build -t docker.ghiasi.me:5000/go-cloud:0.1-2 -f deploy/Dockerfile .

RUN echo "Asia/Tehran" > /etc/timezone-
RUN dpkg-reconfigure -f noninteractive tzdata

RUN useradd -ms /bin/bash newuser

USER newuser
WORKDIR /home/newuser

ARG docker_version=0.0-1

COPY ./main /home/newuser/main

CMD ["/home/newuser/main"]