FROM alpine:latest
WORKDIR /app
COPY ./build /app/
EXPOSE 80
Run the Go application.
CMD ["bash","/app/build"]
