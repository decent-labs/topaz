FROM golang:1.11.1

ENV topaz_root=/go/src/github.com/decentorganization/topaz

WORKDIR ${topaz_root}/models
COPY ./models .

WORKDIR ${topaz_root}/api
COPY ./api .

RUN go get -v
RUN CGO_ENABLED=1 GOOS=linux go build -o api -a -installsuffix cgo

WORKDIR ${topaz_root}/api/settings
COPY ./api/settings/pre.json .
COPY ./api/settings/prod.json .
COPY ./api/settings/tests.json .
COPY ./api/settings/keys ./keys

WORKDIR ${topaz_root}/api
CMD ["./api"]
