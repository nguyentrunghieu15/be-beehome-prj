FROM golang:1.22 as builder

WORKDIR /usr/src/app

COPY . /usr/src/app/
RUN go mod tidy
RUN go build -v -o /opt/app /usr/src/app/pro-manager-service/cmd


FROM alpine:3.14
RUN apk add libc6-compat
COPY --from=builder /opt/app /app
RUN mkdir -p /pro-manager-service/logs
ENTRYPOINT ["/app"]