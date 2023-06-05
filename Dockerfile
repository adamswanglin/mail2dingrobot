FROM golang:1.20.4-alpine3.18 as builder

# 编辑区域
ARG HOME=/app
ARG EXECUTABLE=mail2dingrobot

ENV GOPROXY=https://goproxy.cn,direct
WORKDIR $HOME
COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $EXECUTABLE .

FROM alpine:latest
ARG USER=appuser
ARG USER_ID=2001
ARG GROUP=${USER}
ARG GROUP_ID=${USER_ID}
ARG EXECUTABLE=mail2dingrobot

RUN addgroup -g ${GROUP_ID} -S ${GROUP} && adduser -S -D -G ${GROUP} -u ${USER_ID} ${USER} -h $HOME
WORKDIR app
COPY --from=builder app/mail2dingrobot mail2dingrobot
COPY --from=builder app/config.yml config.yml

RUN chown $USER:$GROUP /app/* && chmod +x /app/mail2dingrobot
USER $USER
CMD ["/app/mail2dingrobot", "--alsologtostderr", "-c", "/app/config.yml"]
EXPOSE 1025