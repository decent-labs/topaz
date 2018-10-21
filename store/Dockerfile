FROM golang:1.11.1 as builder

WORKDIR /go/src/github.com/decentorganization/topaz/store

COPY ./store .
COPY ./models/main.go ./models.go

RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build  -o store -a -installsuffix cgo main.go models.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/decentorganization/topaz/store/store .

CMD ["./store"]
