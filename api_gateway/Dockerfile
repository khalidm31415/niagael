FROM golang:1.18-alpine AS build

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/api_gateway

FROM alpine:latest
RUN apk add bash
COPY --from=build /bin/api_gateway /bin/api_gateway
COPY .env .
CMD ["/bin/api_gateway"]
