# primeiro builda o projeto
FROM golang:1.22-alpine AS build

WORKDIR /usr/src/app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o server

# depois roda
FROM alpine:3

WORKDIR /usr/src/app

COPY --from=build /usr/src/app .

CMD ["/usr/src/app/server"]

# verificar necessidade
EXPOSE 9000