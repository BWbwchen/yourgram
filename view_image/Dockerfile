FROM golang:1.16 AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /view_image  ./cmd/main.go

CMD [ "go", "test", "-v", "-cover", "yourgram/view_image/service"]

FROM alpine:latest AS deploy

WORKDIR /
COPY --from=build /view_image /
CMD ["/view_image"]