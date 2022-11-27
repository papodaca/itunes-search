FROM golang:1.19 as build

WORKDIR /usr/local/go/src/app
COPY . .

RUN go install github.com/GeertJohan/go.rice/rice@latest && \
  go mod download && \
  go generate && \
  CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /
ENV GIN_MODE=release
ENTRYPOINT [ "/app" ]
