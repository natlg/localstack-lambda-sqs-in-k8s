FROM alpine:3.7

RUN apk add --no-cache \
        bash \
        bash-completion \
        curl \
        iputils

ADD publisher /usr/local/bin/

EXPOSE 8085
ENTRYPOINT ["publisher"]
CMD [""]