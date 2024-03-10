# primeiro builda o projeto
FROM golang:1.21-alpine AS build

WORKDIR /usr/app

COPY src/go.* ./
RUN go mod download

COPY src/* ./

RUN go build -v -o server

# depois roda
FROM alpine:3

WORKDIR /usr/app

COPY --from=build /usr/app .

CMD ["/usr/app/server"]

# verificar necessidade
EXPOSE 9000