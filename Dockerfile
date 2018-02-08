FROM golang:1.9 as builder

WORKDIR /go/src/github.com/gopenguin/minimal-ldap-proxy

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o minimal-ldap-proxy .


FROM scratch
COPY --from=builder /go/src/github.com/gopenguin/minimal-ldap-proxy/minimal-ldap-proxy .
CMD ["/minimal-ldap-proxy"]
