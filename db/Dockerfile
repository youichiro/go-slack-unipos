FROM golang

RUN apt update && apt install git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /opt/db

CMD echo "hello migrate"
