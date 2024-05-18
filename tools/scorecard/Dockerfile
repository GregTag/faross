FROM golang:1.22.2-alpine3.19
WORKDIR /app
COPY scorecard.go ./
COPY go.mod ./
RUN go build -o /scorecard
ENTRYPOINT ["/scorecard"]