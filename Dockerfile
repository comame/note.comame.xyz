FROM ubuntu:22.04

WORKDIR /root
COPY ./out out

ENTRYPOINT /root/out/server
