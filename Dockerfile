FROM docker.io/library/golang:1.19-buster AS builder

ENV TZ=Asia/Shanghai
#ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY . .

RUN make build

# Bin
# FROM scratch AS bin
# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
FROM docker.io/library/alpine AS bin

ENV TZ=Asia/Shanghai
#RUN apk add --no-cache fish

COPY --from=builder /src/data /var/lib/gosd/data
COPY --from=builder /src/gosd /usr/bin/gosd

#RUN gosd completion fish > /etc/fish/completions/gosd.fish

EXPOSE 8000/tcp

WORKDIR /var/lib/gosd

ENTRYPOINT ["/usr/bin/gosd"]

CMD ["server"]
