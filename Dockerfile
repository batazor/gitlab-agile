FROM golang:1.12-alpine as builder

ENV GO111MODULE on

# Build project
WORKDIR /go/src/github.com/batazor/gitlab-agile
COPY . .
RUN apk add --update git && \
  go get -u github.com/gobuffalo/packr/packr && \
  packr build cmd/gitlab-agile/main.go && \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build \
  -a \
  -mod vendor \
  -installsuffix cgo \
  -o gitlab-agile ./cmd/gitlab-agile

FROM alpine

USER 10001

WORKDIR /app/
COPY --from=builder /go/src/github.com/batazor/gitlab-agile/gitlab-agile .
ENTRYPOINT ["./gitlab-agile"]
