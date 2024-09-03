FROM ubuntu:22.04

WORKDIR /root
COPY ./out out
COPY ./static static
COPY ./templates templates

ENTRYPOINT /root/out/server
