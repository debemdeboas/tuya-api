FROM golang:1.22-bookworm

ENV DEBIAN_FRONTEND=noninteractive

EXPOSE 8015

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /tuya-api
CMD ["/tuya-api"]