FROM alpine
RUN   sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'  /etc/apk/repositories && apk update && apk add --no-cache ca-certificates wget
ADD ./certs /certs
ADD ./bin/repimage /repimage
ADD ./scripts/download-allowlist.sh /tmp/download-allowlist.sh
RUN /tmp/download-allowlist.sh && rm /tmp/download-allowlist.sh
ENV ALLOWLIST_FILE=/etc/repimage/allowlist.txt
ENV ALLOWLIST_UPDATE_INTERVAL=1h