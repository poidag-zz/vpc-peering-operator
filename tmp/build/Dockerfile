FROM alpine:3.6

RUN apk --no-cache add ca-certificates
RUN adduser -D vpc-peering-operator
USER vpc-peering-operator

ADD tmp/_output/bin/vpc-peering-operator /usr/local/bin/vpc-peering-operator
