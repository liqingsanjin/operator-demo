FROM ubuntu:16.04

RUN apt-get update

ENV INCLUSTER=yes
ENV	QiniuNginxKind=QiniuNginx
ENV	QiniuNginxPlural=qiniunginxs
ENV	QiniuNginxGroup=example.com
ENV	QiniuNginxVersion=v1
ENV	QiniuNginxName="$QiniuNginxPlural.$QiniuNginxGroup"

COPY ./app /app

ENTRYPOINT /app