FROM golang:1.10-alpine as builder

RUN apk add --no-cache git build-base

WORKDIR /go/src/github.com/gopenguin/minimal-ldap-proxy

RUN go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -tags "json1 fts5 sqlite_omit_load_extension" -a -installsuffix cgo -ldflags "-linkmode external -extldflags \"-static -lc\" -w -s" -o minimal-ldap-proxy .


FROM scratch
COPY --from=builder /go/src/github.com/gopenguin/minimal-ldap-proxy/minimal-ldap-proxy .
CMD ["/minimal-ldap-proxy"]
