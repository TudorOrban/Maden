FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -o myapp

FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/myapp /myapp

CMD ["/myapp"]