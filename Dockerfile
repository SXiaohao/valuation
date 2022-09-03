#FROM golang:1.16 AS builder
#
#COPY . /src
#WORKDIR /src
#
#RUN GOPROXY=https://goproxy.cn make build
#
#FROM debian:stable-slim
#
#RUN apt-get update && apt-get install -y --no-install-recommends \
#		ca-certificates  \
#        netbase \
#        && rm -rf /var/lib/apt/lists/ \
#        && apt-get autoremove -y && apt-get autoclean -y
#
#COPY --from=builder /src/bin /app
#COPY --from=builder /src/configs/config.yaml /app
#WORKDIR /app
#
#EXPOSE 8000
#EXPOSE 9000
#
#CMD ["./valuation", "-conf", "./config.yaml"]
FROM golang:alpine3.15 AS builder

COPY . /src
WORKDIR /src
RUN apk add gcc g++ make cmake gfortran libffi-dev openssl-dev libtool
RUN GOPROXY=https://goproxy.cn make build

FROM alpine:latest
COPY --from=builder /src/bin /app
COPY --from=builder /src/configs/config.yaml /app
WORKDIR /app


CMD ["./valuation", "-conf", "./config.yaml"]