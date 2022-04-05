FROM golang:1.18-alpine AS build

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/product

FROM alpine:latest
RUN apk add bash
COPY --from=build /bin/product /bin/product
COPY .env .
CMD ["/bin/product"]
