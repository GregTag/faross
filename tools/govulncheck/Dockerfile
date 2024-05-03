FROM golang:1.22.2-alpine3.19
WORKDIR /app
COPY vuln-example .
RUN go install golang.org/x/vuln/cmd/govulncheck@v1.1.0
ENTRYPOINT [ "govulncheck", "-json",  "./..."]