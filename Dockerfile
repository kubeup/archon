FROM alpine
RUN apk add --update ca-certificates
ADD archon-controller /archon-controller
