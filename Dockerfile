FROM golang:1.15 AS builder
COPY . /src
WORKDIR /src
RUN make build

FROM debian

RUN apt-get update && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

COPY --from=builder /src/dist/proxy-server /proxy-server
CMD ["/proxy-server"]
