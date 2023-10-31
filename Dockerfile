FROM alpine:latest
WORKDIR /app
COPY ./build /app/
EXPOSE 8080
ENTRYPOINT ["/bin/bash","/app/build"]
