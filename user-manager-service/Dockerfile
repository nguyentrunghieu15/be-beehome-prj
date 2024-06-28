FROM golang:1.22 as builder

WORKDIR /usr/src/app

COPY . /usr/src/app/
RUN go mod tidy
RUN go build -v -o /opt/app /usr/src/app/user-manager-service/cmd


FROM alpine:3.14
COPY --from=builder /opt/app /app
RUN apk add libc6-compat
RUN mkdir -p /user-manager-service/logs
ENTRYPOINT ["/app"]