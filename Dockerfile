FROM gcr.io/moonrhythm-containers/alpine

RUN mkdir -p /app
WORKDIR /app

COPY tunnel-http-socks5 ./

ENTRYPOINT ["/app/tunnel-http-socks5"]
