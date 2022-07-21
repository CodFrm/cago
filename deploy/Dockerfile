FROM golang:1.18-alpine3.16 as build

ARG APP_NAME=cago
ARG GOOS=linux
ARG GOARCH=amd64
ARG GOPROXY="https://goproxy.cn,direct"

WORKDIR /app

COPY . .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
	apk add git

RUN go build -o ${APP_NAME} .

FROM alpine:3.16

ARG APP_NAME=cago

WORKDIR /app

COPY --from=build /app/${APP_NAME} .

RUN ls -l && chmod +x ${APP_NAME}

CMD ["./${APP_NAME}"]
