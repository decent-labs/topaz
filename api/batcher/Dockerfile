FROM golang:1.11.1 as builder

WORKDIR /go/src/github.com/decentorganization/topaz/api/batcher
COPY ./api/batcher .

# Temporary
COPY ./api/core/database ./database
COPY ./api/models ./models

RUN ls -a

RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o batcher -a -installsuffix cgo

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/decentorganization/topaz/api/batcher/batcher .

CMD ["./batcher"]
