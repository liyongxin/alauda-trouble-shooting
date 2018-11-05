FROM golang:1.10 as builder

WORKDIR $GOPATH/src/alauda-trouble-shooting
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o alauda_trouble_shooting .
RUN go install -v


FROM index.alauda.cn/alaudaorg/alaudabase:alpine-supervisor-migrate-1

WORKDIR /

COPY --from=builder /go/src/alauda-trouble-shooting .

RUN chmod +x /alauda_trouble_shooting
EXPOSE 3322

CMD ["/alauda_trouble_shooting"]
