FROM moodysan/gobaseimage:latest
WORKDIR /app
COPY ./build /app/
EXPOSE 80
Run the Go application.
CMD ["/app/build"]