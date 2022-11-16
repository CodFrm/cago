FROM alpine:3.16

ARG APP_NAME=cago

WORKDIR /app

COPY ${APP_NAME} .

RUN ls -l && chmod +x ${APP_NAME}

CMD ["./${APP_NAME}"]
