FROM golang:latest

WORKDIR /home/go/app

# Installing goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD tail -f /dev/null