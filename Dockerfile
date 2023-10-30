FROM golang:1.18
WORKDIR /app
COPY ./build /app/
EXPOSE 8080
CMD ["/app/build"]
