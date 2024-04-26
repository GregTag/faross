FROM golang:1.22.2-alpine3.19 AS build

RUN apk add --no-cache gcc g++

WORKDIR /app

COPY . .

WORKDIR /app/firewall

RUN CGO_ENABLED=1 GOOS=linux go build -o /firewall ./cmd/app/main.go

FROM docker:26.1.0-cli AS runtime

WORKDIR /app

COPY firewall/config config
COPY --from=build /firewall firewall

ENTRYPOINT [ "./firewall" ]
