FROM golang:1.14 as go-build

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN GO111MODULE=on GOPROXY=https://goproxy.io CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main *.go


FROM registry.cn-shenzhen.aliyuncs.com/happysooner/golang-runtime

WORKDIR /

COPY --from=go-build /app/main /

ENTRYPOINT ["/main"]