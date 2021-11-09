FROM golang:1.14.6-alpine3.12 as builder

COPY go.mod go.sum /go/src/telegramStravaBot/
WORKDIR /go/src/telegramStravaBot
RUN go mod download
COPY . /go/src/telegramStravaBot
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/telegramStravaBot telegramStravaBot

RUN apk add --no-cache bash git openssh


FROM alpine

RUN apk add --no-cache  make musl-dev go bash git openssh ca-certificates && update-ca-certificates
COPY --from=builder /go/src/telegramStravaBot/build/telegramStravaBot /usr/bin/telegramStravaBot

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Glide

WORKDIR $GOPATH

CMD ["make"]

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/telegramStravaBot"]
