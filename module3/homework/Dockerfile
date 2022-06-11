FROM golang:1.18-alpine AS build
COPY ./project /go/project
WORKDIR /go/project
RUN go build -o /bin/project

FROM alpine:latest as prod
COPY --from=build /bin/project .
EXPOSE 8080
