FROM golang:1.22 as builder

WORKDIR /usr/src/app

COPY . /usr/src/app/
RUN go mod tidy
RUN go build -v -o /opt/app /usr/src/app/authorize-service/cmd


FROM alpine:3.14
RUN apk add libc6-compat
COPY --from=builder /opt/app /app
RUN mkdir -p /authorize-service/logs
ENTRYPOINT ["/app"]