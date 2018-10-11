FROM golang:1.11.1

WORKDIR /go/src/github.com/decentorganization/topaz/ethereum

COPY . .

RUN go get -v
RUN CGO_ENABLED=1 GOOS=linux go build -o ethereum -a -installsuffix cgo main.go

CMD ["./ethereum"]
