FROM surnet/alpine-wkhtmltopdf:3.9-0.12.5-full as wkhtmltopdf
FROM golang:1.14.2-alpine3.11

LABEL MAINTAINER="janeysesions@gmail.com"

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/,direct" \
  GO111MODULE=on

WORKDIR /go/src/github.com/jackdon/gowxapi
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o app .

# Final stage:
FROM golang:1.14.2-alpine3.11

WORKDIR /root

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# Install dependencies for wkhtmltopdf
RUN apk add --no-cache \
  libstdc++ \
  libx11 \
  libxrender \
  libxext \
  libssl1.1 \
  ca-certificates \
  fontconfig \
  freetype \
  ttf-dejavu \
  ttf-droid \
  ttf-freefont \
  ttf-liberation \
  ttf-ubuntu-font-family \
&& apk add --no-cache --virtual .build-deps \
  msttcorefonts-installer \
\
# Install microsoft fonts
&& update-ms-fonts \
&& fc-cache -f \
\
# Clean up when done
&& rm -rf /tmp/* \
&& apk del .build-deps 

# 创建目录,保存代码
RUN mkdir -p /data/gowxapi/logs \
  && mkdir -p /opt/conf/

COPY --from=1 /go/src/github.com/jackdon/gowxapi/app .

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=wkhtmltopdf /bin/wkhtmltoimage /bin/wkhtmltoimage
COPY --from=wkhtmltopdf /bin/libwkhtmltox* /bin/

# 配置文件
ENV CONFIG=/opt/conf/config.yaml

VOLUME [ "/opt/conf", "/opt/lesson", "/data/gowxapi/logs" ]  

COPY config.yaml /opt/conf/config.yaml

EXPOSE 8080

CMD ["./app"]