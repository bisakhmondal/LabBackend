FROM golang:alpine
LABEL maintainer="Bisakh Mondal <bisakhmondal00@gmail.com>"

RUN apk add \
        bash \
        curl \
        zip

COPY ./auth /backend
WORKDIR /backend/

ENV CertKey certs/server.key
ENV CertFile certs/server.crt
ENV API_SECRET O1z378nx3Nu3o2Hf0DeYTIwBBJIuEPJHrYmbLf2wRJijK7v5oO

EXPOSE 8080

CMD ["go","run","main.go"]
