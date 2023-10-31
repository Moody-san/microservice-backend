FROM alpine:latest
WORKDIR /app
COPY ./build /app/
EXPOSE 8080
CMD ["/app/build"]
