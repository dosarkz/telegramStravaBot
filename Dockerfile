FROM golang:1.16-alpine3.14 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -o /build


FROM alpine:3.14

ARG USER=default
ENV HOME /home/$USER

# install sudo as root
RUN apk add --update sudo

# add new user
RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER

WORKDIR /

COPY --from=build /build /build
COPY --from=build  /app/.env /.env
COPY --from=build  /app/wait-for-postgres.sh /wait-for-postgres.sh
COPY ./data/database/migrations/ /data/database/migrations/

EXPOSE 8000

ENTRYPOINT ["/build"]