FROM alpine:latest
RUN mkdir /app

COPY broker-service/brokerServiceApp /app

# Run the server executable
CMD [ "/app/brokerServiceApp" ]