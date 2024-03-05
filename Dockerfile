FROM golang:1.22.0-alpine3.18 AS build-stage
WORKDIR /app
COPY ./ /app
RUN mkdir -p /app/build
RUN go mod download
RUN go build -v -o /app/build/api ./cmd/api
EXPOSE 8080
CMD ["/app/build/api"]
