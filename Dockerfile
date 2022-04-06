FROM golang:1.17-buster AS builder

ENV TZ=Asia/Shanghai
#ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY . .

RUN make build

# Bin
# FROM scratch AS bin
# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
FROM alpine AS bin

COPY --from=builder /src/data /var/lib/gosd/data
COPY --from=builder /src/gosd /usr/bin/gosd

RUN apk add --no-cache bash
COPY --from=builder /go/pkg/mod/github.com/urfave/cli/v2@v2.3.0/autocomplete/bash_autocomplete /etc/bash/autocomplete.d/gosd
RUN printf 'export PS1="\u@\h:\w\$ " && PROG=gosd source /etc/bash/autocomplete.d/gosd' >> /root/.bashrc

ENV TZ=Asia/Shanghai

EXPOSE 8000/tcp

WORKDIR /var/lib/gosd

ENTRYPOINT ["/usr/bin/gosd"]
