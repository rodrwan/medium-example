FROM alpine

WORKDIR /app

COPY bin/medium-example /usr/bin/medium-example

ENTRYPOINT [ "medium-example" ]
CMD [ " -postgres-dsn=$POSTGRES_DSN" ]
