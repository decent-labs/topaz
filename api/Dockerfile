FROM golang:1.11.1 as builder

ENV topaz_root=/go/src/github.com/decentorganization/topaz

WORKDIR ${topaz_root}/models
COPY ./models .

WORKDIR ${topaz_root}/api
COPY ./api .

RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o topazapi -a -installsuffix cgo

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /go/src/github.com/decentorganization/topaz/api/topazapi .

WORKDIR /app/settings
COPY ./api/settings/pre.json .
COPY ./api/settings/prod.json .
COPY ./api/settings/tests.json .
COPY ./api/settings/keys ./keys

WORKDIR /app
CMD ["./topazapi"]
