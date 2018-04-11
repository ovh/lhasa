FROM debian:jessie

COPY dist /app

WORKDIR /app

RUN mkdir -p /etc/appcatalog
RUN chown -R nobody:nogroup /app && chmod +x /app/appcatalog-configuration
RUN chown -R nobody:nogroup /app && chmod +x /app/appcatalog
USER nobody

EXPOSE 8081

CMD ./appcatalog-configuration --debug --output=./config.json start && ./appcatalog --config=./config.json
