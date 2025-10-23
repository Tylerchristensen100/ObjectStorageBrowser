FROM golang:1.25.1-alpine AS build

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates


WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN go build -o server .




FROM scratch


COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /app/etc/ /etc
COPY --from=build /app/server .
EXPOSE 3000
CMD ["./server"]