FROM golang:1.22.3-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o bin/app cmd/app/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/app .

CMD ["sh", "-c", "./app < ./file", "sh"]
