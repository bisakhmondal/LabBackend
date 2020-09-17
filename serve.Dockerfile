FROM golang:alpine
LABEL maintainer="Bisakh Mondal <bisakhmondal00@gmail.com>"

RUN apk add \
        bash \
        curl \
        zip

COPY ./serving-api /backend
WORKDIR /backend/

ENV CertKey certs/server.key
ENV CertFile certs/server.crt

EXPOSE 8080

CMD ["go","run","main.go"]
