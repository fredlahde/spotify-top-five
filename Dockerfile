FROM ubuntu:latest

ADD spotify-top-five /usr/bin/spotify-top-five

ENTRYPOINT /usr/bin/spotify-top-five