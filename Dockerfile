FROM golang:1.20-alpine

COPY . /go/src/app

WORKDIR /go/src/app/cmd/audit-log

RUN go build -o app main.go

EXPOSE 8080

CMD ["./app"]