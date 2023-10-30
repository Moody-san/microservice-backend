FROM golang:1.18
WORKDIR /app
COPY ./build /app/
EXPOSE 8080
Run the Go application.
CMD ["/app/build"]
