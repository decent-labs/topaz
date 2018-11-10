FROM golang:1.11.1

ENV topaz_root=/go/src/github.com/decentorganization/topaz

WORKDIR ${topaz_root}/api
COPY ./api .

RUN go get -v
RUN CGO_ENABLED=1 GOOS=linux go build -o api -a -installsuffix cgo

COPY ./api/settings/pre.json ./settings
COPY ./api/settings/prod.json ./settings
COPY ./api/settings/tests.json ./settings
COPY ./api/settings/keys ./settings/keys

CMD ["./api"]
