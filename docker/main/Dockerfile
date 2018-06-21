FROM debian:jessie

COPY dist /app

WORKDIR /app

RUN ls -lrt /app
RUN mkdir -p /etc/appcatalog
RUN chown -R nobody:nogroup /app && chmod +x /app/mycompany.sh
RUN chown -R nobody:nogroup /app && chmod +x /app/appcatalog.sh
RUN chown -R nobody:nogroup /app && chmod +x /app/appcatalog-configuration
RUN chown -R nobody:nogroup /app && chmod +x /app/appcatalog
USER nobody

EXPOSE 8081

CMD ./appcatalog-configuration --debug --output=./config.json start && ./appcatalog.sh --config=./config.json
