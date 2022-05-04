FROM golang:1.18-alpine3.15 AS  build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN --mount=type=cache,target=/var/.cache/apk \
    apk update && apk add --no-cache bash git openssh build-base

COPY . .


RUN GOOS=linux go build -o /build


FROM alpine:3.15

RUN apk update && apk add --no-cache make musl-dev go bash git openssh ca-certificates && update-ca-certificates

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR $GOPATH

COPY --from=build /build /build
COPY --from=build  /app/.env .
COPY database/migrations/ $GOPATH/database/migrations

EXPOSE 8000

ENTRYPOINT ["/build"]