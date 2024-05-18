FROM node:22.1.0-alpine3.19 as build-npm

WORKDIR /build

COPY frontend frontend

WORKDIR /build/frontend

RUN npm install
RUN npm run build


FROM golang:1.22.2-alpine3.19 AS build-go

RUN apk add --no-cache gcc g++

WORKDIR /build

COPY . .

WORKDIR /build/firewall

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/app/main.go


FROM docker:26.1.0-cli AS runtime

WORKDIR /app

COPY firewall/config.yaml config.yaml
COPY gatherlaunch/instruments.json instruments.json
COPY --from=build-npm /build/frontend/build static 
COPY --from=build-go /build/firewall/main main


ENTRYPOINT [ "./main" ]
