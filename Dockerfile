FROM golang:1.17-alpine as builder

ADD ./ /web-screenshot

RUN apk update \
  && apk add git gcc make \
  && cd /web-screenshot \
  && make build

FROM alpine

# tzdata 安装所有时区配置或可根据需要只添加所需时区

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates tzdata

COPY --from=builder /web-screenshot/web-screenshot /usr/local/bin/web-screenshot
COPY --from=builder /web-screenshot/entrypoint.sh /entrypoint.sh

EXPOSE 7000

USER go

WORKDIR /home/go

HEALTHCHECK --timeout=10s --interval=10s CMD [ "wget", "http://127.0.0.1:7000/ping", "-q", "-O", "-"]

CMD ["web-screenshot"]

ENTRYPOINT ["/entrypoint.sh"]