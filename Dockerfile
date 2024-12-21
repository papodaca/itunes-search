FROM golang:1.23 AS build

WORKDIR /go/src/itunes-search
COPY . .

RUN go install .

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/* /
ENV GIN_MODE=release
ENTRYPOINT [ "/itunes-search" ]
