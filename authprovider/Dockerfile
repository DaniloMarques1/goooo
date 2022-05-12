FROM golang:1.18

WORKDIR ./authProvider

COPY . .

RUN go mod download

RUN go build .

CMD ["./authprovider"]
