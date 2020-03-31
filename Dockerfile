FROM golang:1.14-alpine as builder
WORKDIR /src/github.com/codebender/qrcode-api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o qrcode-api main.go

FROM alpine:3.11
RUN adduser -S -D -H -h /app appuser
USER appuser
WORKDIR /app
COPY --from=builder /src/github.com/codebender/qrcode-api/qrcode-api .
CMD ["./qrcode-api"]