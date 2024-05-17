FROM golang:1.22.3-alpine

WORKDIR /app

COPY . .

RUN go build -o bin/app cmd/app/main.go

CMD ["sh", "-c", "./bin/app < ./data/file"]
