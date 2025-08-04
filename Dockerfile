FROM golang:1.24

RUN mkdir app

WORKDIR /app

ADD . /app

RUN go build -o main cmd/main.go

EXPOSE 8000

CMD ["./main"]