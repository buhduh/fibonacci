FROM golang:1.16-buster AS reserve_trust_build
ARG LD_FLAGS
ARG PACKAGE
WORKDIR /go/src
COPY ./src .
RUN go build -ldflags="$LD_FLAGS" -o /$PACKAGE $PACKAGE

FROM golang:1.16-buster AS reserve_trust_test
ARG LD_FLAGS
ARG PACKAGE
WORKDIR /go/src
COPY ./src .
RUN go test -ldflags="$LD_FLAGS" -v $PACKAGE

#probably silly, just felt like separating the application containter
FROM debian:buster AS reserve_trust_application
ARG PACKAGE
COPY --from=reserve_trust_build /$PACKAGE /$PACKAGE
ENTRYPOINT ["/fibonacci"]
