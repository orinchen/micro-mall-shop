FROM golang:1.16.15-buster as gobuild
WORKDIR $GOPATH/src/gitee.com/cristiane/micro-mall-shop
COPY . .
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on
RUN bash ./build.sh
# FROM alpine:latest as gorun
FROM centos:latest as gorun
WORKDIR /www/
COPY --from=gobuild /go/src/gitee.com/cristiane/micro-mall-shop/micro-mall-shop .
COPY --from=gobuild /go/src/gitee.com/cristiane/micro-mall-shop/etc ./etc
CMD ["./micro-mall-shop"]
