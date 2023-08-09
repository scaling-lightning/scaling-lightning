# Run from the root of scaling lightning
# docker build -t <user>/lnd-client -f clients/lnd/Dockerfile .

# build stage
FROM golang:buster as backend-builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY pkg/ ./pkg
COPY clients/lnd/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lnd-client .

# final stage
FROM debian:buster-slim
COPY --from=backend-builder /app/lnd-client /app/
RUN useradd -ms /bin/bash lnd
RUN apt-get -y update && apt-get -y --no-install-recommends install ca-certificates bash && rm -rf /var/lib/apt/lists/*;
RUN update-ca-certificates
WORKDIR /app
USER lnd
ENTRYPOINT ["./lnd-client"]
