FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY mail-service/mailServiceApp /app
COPY mail-service/templates/. /templates

# Run the server executable
CMD [ "/app/mailServiceApp" ]