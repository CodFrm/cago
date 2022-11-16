FROM alpine:3.16

ARG APP_NAME=cago

ENV APP_NAME=$APP_NAME

WORKDIR /app

COPY $APP_NAME .

RUN ls -l && chmod +x $APP_NAME

CMD ["sh", "-c", "./$APP_NAME"]
