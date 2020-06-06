FROM debian:9.4-slim

RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y python-pip
RUN pip install awscli-local
RUN apt-get install -y curl

ADD scripts scripts/
ADD worker.zip worker.zip
