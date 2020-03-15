FROM scratch
VOLUME /tmp/apps/logs
COPY ./amqpserver /tmp/apps/amqpserver
COPY ./config.json /tmp/apps/config.json
WORKDIR /tmp/apps
CMD ["./amqpserver","-c","/tmp/apps/config.json"]
