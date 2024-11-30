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

RUN if [ "$USE_PROFILE" == "prod" ]; then \
      cp /build/.env.prod /app/.env \
      cp /build/firebase.json /app/ \
      echo "use prod" \
    elif [ "$USE_PROFILE" == "live" ]; then \
      cp /build/.env.live /app/.env \
      cp /build/test-firebase.json /app/ \
      echo "use live" \
    fi

EXPOSE 8080
ENV BREEZENOTE_PROFILE=${USE_PROFILE}

CMD ["/bin/breezenote"]