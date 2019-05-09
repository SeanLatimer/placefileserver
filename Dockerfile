FROM golang:1.12-alpine as builder

WORKDIR /placefileserver
COPY . .
RUN apk add --no-cache git make
RUN go get -u github.com/gobuffalo/packr/v2/packr2
RUN make clean && make build

FROM alpine:latest
WORKDIR /app/
RUN apk add --no-cache ca-certificates
COPY --from=builder /placefileserver/build .
CMD [ "./placefileserver" ]
