FROM golang:1.16-buster AS builder

ENV TZ=Asia/Shanghai
#ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY . .

RUN make build

# Bin
FROM scratch AS bin

COPY --from=builder /src/gosd /usr/bin/gosd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Shanghai

EXPOSE 8000/tcp

WORKDIR /var/lib/gosd

ENTRYPOINT ["/usr/bin/gosd", "container"]
