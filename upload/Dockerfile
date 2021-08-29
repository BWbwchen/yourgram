FROM golang:1.16 AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /upload  ./cmd/main.go

CMD [ "go", "test", "-v", "-cover", "yourgram/upload/service"]

FROM alpine:latest AS deploy

WORKDIR /
COPY --from=build /upload /
CMD ["/upload"]