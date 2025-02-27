FROM golang:1.16.5-alpine3.12 as builder

WORKDIR /workdir

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /go-app server.go

FROM alpine:3.21
COPY --from=builder /go-app .
ENTRYPOINT ["./go-app"]
