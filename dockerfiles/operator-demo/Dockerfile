FROM ubuntu:16.04

RUN apt-get update
RUN go build -o ./app ~/go/src/operator-demo/operator.go

COPY ./app /app

ENTRYPOINT /app
