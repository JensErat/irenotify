# builder image
FROM golang:1.20 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o irenotify .

FROM scratch
COPY --from=builder /build/irenotify .
ENTRYPOINT [ "/irenotify" ]