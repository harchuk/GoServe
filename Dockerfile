FROM golang:1.20-alpine as builder
LABEL maintainer="Harchuk <harchuk1@gmail.com>"
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY UserServer.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main UserServer.go
FROM alpine:3
COPY --from=builder main /bin/main
EXPOSE 8080
ENTRYPOINT ["/bin/main"]