FROM golang:1.17 AS build
COPY ./project /go/project
WORKDIR /go/project
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN GOOS=linux go build -installsuffix cgo -o /bin/project

FROM busybox
COPY --from=build /bin/project /usr/local/httpserver/
COPY ./project/config/config.yml /usr/local/httpserver/config/config.yml
EXPOSE 8080
ENV ENV local
WORKDIR /usr/local/httpserver/
ENTRYPOINT ["/usr/local/httpserver/project"]
