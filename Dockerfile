FROM golang:1.21.4-alpine3.18
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app1

FROM alpine:latest
RUN apk --no-cache add bash
WORKDIR /app
COPY --from=builder /app1 /app1
EXPOSE 8080
CMD ["/app1"]
