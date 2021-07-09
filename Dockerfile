FROM golang:1.16-buster AS builder

WORKDIR /src

COPY . .

RUN make build

# Bin
FROM alpine AS bin

COPY --from=builder /src/gosd /usr/bin/gosd

EXPOSE 8000/tcp

ENTRYPOINT ["/usr/bin/gosd"]
