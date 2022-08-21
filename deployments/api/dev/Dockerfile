FROM golang:1.18.1-alpine

RUN apk --no-cache add build-base git gcc openssh-client curl

WORKDIR /app

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air -c .air.toml
