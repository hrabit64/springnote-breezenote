FROM golang:1.23
RUN apt-get -y update && apt-get install -y bash libwebp-dev


WORKDIR /build
COPY . .
RUN go build -o /bin/breezenote ./cmd/main.go

RUN chmod +x /bin/breezenote

WORKDIR /app
COPY .env.prod .
COPY firebase.json .

RUN mkdir ./logs
RUN mkdir ./data
RUN mkdir ./data/images


EXPOSE 8080
ENV BREEZENOTE_PROFILE="prod"

CMD ["/bin/breezenote"]