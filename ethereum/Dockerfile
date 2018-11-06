FROM golang:1.11.1


ENV topaz_root=/go/src/github.com/decentorganization/topaz

WORKDIR ${topaz_root}/models
COPY ./models .

WORKDIR ${topaz_root}/ethereum
COPY ./ethereum .

RUN go get -v
RUN CGO_ENABLED=1 GOOS=linux go build -o ethereum -a -installsuffix cgo

CMD ["./ethereum"]
