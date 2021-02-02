# builder image
FROM golang:1.15-alpine3.13 as builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -o mcserver -ldflags="-s -w" ./cmd/mcserver


# generate clean, final image for end users
FROM alpine:3.13
MAINTAINER tim.satke@gmail.com
EXPOSE 25565

COPY --from=builder /build/mcserver .

# executable
ENTRYPOINT [ "./mcserver" ]
# arguments that can be overridden
CMD [ ":25565" ]