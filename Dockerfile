FROM gcr.io/moonrhythm-containers/alpine

RUN mkdir -p /app
WORKDIR /app

COPY tunnel-http-socks5 ./
COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

CMD ["sh", "/app/entrypoint.sh"]
