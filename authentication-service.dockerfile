FROM alpine:latest
RUN mkdir /app

COPY authentication-service/authServiceApp /app

# Run the server executable
CMD [ "/app/authServiceApp" ]