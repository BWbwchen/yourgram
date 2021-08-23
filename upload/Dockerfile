FROM golang:1.16 AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /upload 

CMD [ "go", "test", "-v", "yourgram/upload/service"]

FROM alpine:latest AS deploy

WORKDIR /
COPY --from=build /upload /
CMD ["/upload"]