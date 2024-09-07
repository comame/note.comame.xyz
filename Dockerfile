FROM ubuntu:22.04

RUN apt update -y && apt install ca-certificates -y

WORKDIR /root
COPY ./out out
COPY ./static static
COPY ./templates templates

ENTRYPOINT /root/out/server
