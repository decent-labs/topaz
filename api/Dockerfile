FROM golang:1.11.1 as builder

ENV topaz_root=/go/src/github.com/decentorganization/topaz

WORKDIR ${topaz_root}/models
COPY ./models .

WORKDIR ${topaz_root}/api
COPY ./api .

RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build  -o api -a -installsuffix cgo

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/decentorganization/topaz/api/api .

CMD ["./api"]
