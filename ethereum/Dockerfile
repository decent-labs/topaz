FROM golang:1.11.1 as builder

WORKDIR /go/src/github.com/decentorganization/topaz/ethereum

COPY . .

RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o ethereum -a -installsuffix cgo main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/decentorganization/topaz/ethereum/ethereum .

CMD ["./ethereum"]
