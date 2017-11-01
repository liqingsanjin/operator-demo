FROM ubuntu:16.04

RUN apt-get update

ENV INCLUSTER=yes
	QiniuNginxKind=QiniuNginx
	QiniuNginxPlural=qiniunginxs
	QiniuNginxGroup=example.com
	QiniuNginxVersion=v1
	QiniuNginxName="$QiniuNginxPlural.$QiniuNginxGroup"
	namespace=default

COPY ./app /app

ENTRYPOINT /app