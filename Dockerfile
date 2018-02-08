FROM golang:1.9 as builder

WORKDIR /go/src/github.com/gopenguin/minimal-ldap-proxy

RUN go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o minimal-ldap-proxy .


FROM scratch
COPY --from=builder /go/src/github.com/gopenguin/minimal-ldap-proxy/minimal-ldap-proxy .
CMD ["/minimal-ldap-proxy"]
