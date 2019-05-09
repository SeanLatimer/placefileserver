FROM golang:1.12-alpine as builder

WORKDIR /placefileserver
COPY . .
RUN apk add --no-cache git make
RUN make clean && make build

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/github.com/seanlatimer/placefileserver/* .
CMD [ "./placefileserver" ]
