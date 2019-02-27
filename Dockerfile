FROM golang:1.12-alpine

ADD . /go/src/informator-69-bot

RUN cd /go/src/informator-69-bot && \
    go build -o /srv/informator-69-bot ./app


WORKDIR /srv

CMD ["/srv/informator-69-bot"]

