FROM gcr.io/moonrhythm-containers/alpine

RUN mkdir -p /app
WORKDIR /app

COPY tunnel-http-socks5 ./
COPY entrypoint.sh ./

ENTRYPOINT ["/app/entrypoint.sh"]
