FROM debian:jessie

COPY dist /app

WORKDIR /app

RUN chown -R nobody:nogroup /app \
    && chmod +x /app/mycompany.sh \
                /app/appcatalog.sh \
                /app/appcatalog-configuration \
                /app/appcatalog \
    && apt-get update \
    && apt-get install -y curl wget ca-certificates

USER nobody

EXPOSE 8081

CMD ./appcatalog-configuration --debug --output=./config.json start && ./appcatalog.sh --config=./config.json
