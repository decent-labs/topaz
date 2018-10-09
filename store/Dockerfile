FROM golang:1.11.1 as builder

WORKDIR /go/src/github.com/decentorganization/topaz/store

COPY . .

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build  -o store -a -installsuffix cgo main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/decentorganization/topaz/store/store .

CMD ["./store"]
