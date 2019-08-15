FROM golang:1.12.8-alpine3.10 as builder

RUN apk add --update --no-cache libc6-compat \
        git \
        ca-certificates \
        build-base

RUN mkdir /topaz 
COPY . /topaz/

WORKDIR /topaz/api
RUN go get -u github.com/gobuffalo/packr/packr
RUN packr build -mod=vendor -ldflags "-linkmode external -extldflags -static" -tags netgo -installsuffix netgo -a -o ../topaz-api .

WORKDIR /topaz

FROM alpine:3.10.1
COPY --from=builder /topaz/topaz-api .
RUN apk add --update --no-cache libc6-compat ca-certificates
CMD ["./topaz-api"]
