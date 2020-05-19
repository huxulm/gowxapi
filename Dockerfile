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
FROM alpine

WORKDIR /root

# 创建目录,保存代码
RUN mkdir -p /data/gowxapi/logs \
  && mkdir -p /opt/conf/

COPY --from=0 /go/src/github.com/jackdon/gowxapi/app .

# 配置文件
ENV CONFIG=/opt/conf/config.yaml

VOLUME [ "/opt/conf", "/data/gowxapi/logs" ]  

COPY config.yaml /opt/conf

EXPOSE 8080

CMD ["./app"]