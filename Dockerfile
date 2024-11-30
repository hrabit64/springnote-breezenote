FROM golang:1.23
RUN apt-get -y update && apt-get install -y bash libwebp-dev


WORKDIR /build
COPY . .
RUN go build -o /bin/breezenote ./cmd/main.go

RUN chmod +x /bin/breezenote

WORKDIR /app
ENV USE_PROFILE prod

RUN mkdir ./logs
RUN mkdir ./data
RUN mkdir ./data/images

EXPOSE 8080
ENV BREEZENOTE_PROFILE=${USE_PROFILE}

COPY run_server.sh /run_server.sh
RUN chmod +x /run_server.sh

CMD ["/run_server.sh"]