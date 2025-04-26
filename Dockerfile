FROM docker.m.daocloud.io/alpine:latest

RUN mkdir -p /app

#RUN apk add --no-cache gcompat
#RUN apk add --no-cache tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

COPY ./logto.cleaner /app

WORKDIR /app
ENTRYPOINT ["./logto.cleaner"]
