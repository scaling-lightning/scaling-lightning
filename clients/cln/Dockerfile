# Run from the root of scaling lightning
# docker build -t <user>/cln-client -f clients/cln/Dockerfile .

# build stage
FROM golang:buster as backend-builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY pkg ./pkg
COPY clients/cln ./clients/cln
RUN CGO_ENABLED=0 GOOS=linux go build -o cln-client ./clients/cln

# final stage
FROM debian:buster-slim
COPY --from=backend-builder /app/cln-client /app/
RUN useradd -ms /bin/bash cln
RUN apt-get -y update && apt-get -y --no-install-recommends install ca-certificates bash && rm -rf /var/lib/apt/lists/*;
RUN update-ca-certificates
WORKDIR /app
USER cln
ENTRYPOINT ["./cln-client"]
