FROM golang:1.16-buster AS builder

WORKDIR /src

COPY ./ .

#RUN export GOPROXY=https://goproxy.io,direct make build
RUN make build

# Bin
FROM alpine AS bin

COPY --from=builder /src/gosd /usr/bin/gosd

ENTRYPOINT /usr/bin/gosd
