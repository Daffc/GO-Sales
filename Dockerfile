FROM golang:latest

WORKDIR /home/gouser/app

# Creates `gouser` and `appgroup` and sets permissions on the app`s directory
RUN addgroup appgroup --gid 1000 && \
    useradd gouser --uid 1000 --gid appgroup --home-dir /home/gouser/ && \
    chown -R gouser:appgroup /home/gouser/ 

# All the following commands will be executed by `gouser`, instead of `root`
USER gouser

# Copy artifacts from the build stage and set `gouser` as the owner
COPY --chown=appuser:appgroup . /home/gouser/app

# Installing APP dependencies
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

CMD tail -f /dev/null