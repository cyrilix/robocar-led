FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.21-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /go/src
ADD . .

RUN GOOS=$(echo $TARGETPLATFORM | cut -f1 -d/) && \
    GOARCH=$(echo $TARGETPLATFORM | cut -f2 -d/) && \
    GOARM=$(echo $TARGETPLATFORM | cut -f3 -d/ | sed "s/v//" ) && \
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build -mod vendor -tags netgo,no_d2xx ./cmd/rc-led/


#ARG GOOS=linux
#ARG GOARCH=amd64
#ARG GOARM=""

#RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build -mod vendor -tags netgo ./cmd/rc-led/




FROM gcr.io/distroless/static

USER 1234
COPY --from=builder /go/src/rc-led /go/bin/rc-led
ENTRYPOINT ["/go/bin/rc-led"]
