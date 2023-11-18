FROM --platform=$TARGETPLATFORM golang:1.21.4
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app1
EXPOSE 8080
CMD ["/app1"]
