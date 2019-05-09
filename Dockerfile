FROM golang:1.12-alpine as builder

WORKDIR /placefileserver
COPY . .
RUN apk add --no-cache git make
RUN make clean && make build

FROM alpine:latest
WORKDIR /app/
RUN apk add --no-cache ca-certificates
COPY --from=builder /placefileserver/build .
CMD [ "./placefileserver" ]
