FROM debian:jessie

COPY dist /app

WORKDIR /app

RUN chown -R nobody:nogroup /app && chmod +x /app/appcatalog
USER nobody

EXPOSE 8081

CMD ["./appcatalog"]
