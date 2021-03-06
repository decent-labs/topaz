FROM golang:1.12.8-alpine3.10 as builder

RUN apk add --update --no-cache libc6-compat \
        git \
        ca-certificates \
        build-base

RUN mkdir /topaz 
COPY . /topaz/

WORKDIR /

RUN go get github.com/ethereum/go-ethereum

RUN cd /topaz && \
    cp -r "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" \
          "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"

WORKDIR /topaz/batch

RUN go build -mod=vendor -ldflags "-linkmode external -extldflags -static" -tags netgo -installsuffix netgo -a -o ../topaz-batch .

WORKDIR /topaz/batch

FROM alpine:3.10.1
COPY --from=builder /topaz/topaz-batch .
RUN apk add --update --no-cache libc6-compat ca-certificates
CMD ["./topaz-batch"]  
