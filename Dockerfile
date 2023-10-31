FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /app
COPY ./build /app/
EXPOSE 8080
ENTRYPOINT ["/bin/bash","/app/build"]
