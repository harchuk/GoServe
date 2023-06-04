FROM golang:1.20-alpine as builder
LABEL maintainer="Harchuk <harchuk1@gmail.com>"
WORKDIR /build
COPY go.mod go.sum UserServer.go .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /main UserServer.go
FROM alpine:3
COPY --from=builder main /bin/main
EXPOSE 8080
RUN adduser -D -g '' appuser && chown -R appuser /bin/main
USER appuser
ENTRYPOINT ["/bin/main"]
