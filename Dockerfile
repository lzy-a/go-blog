FROM golang:latest
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE=on
WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .
EXPOSE 8000
CMD ["/app/main"]