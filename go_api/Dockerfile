FROM golang:1.24.0

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
RUN go get github.com/pressly/goose/v3
RUN go get github.com/go-sql-driver/mysql@v1.9.1
RUN go install github.com/swaggo/swag/cmd/swag@latest

CMD tail -f /dev/null