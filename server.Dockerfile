FROM golang:1.19-apline

RUN mkdir /app

add . /app

WORKDIR /app

RUN go build -o main cmd/main.go

CMD ["/app/main"]