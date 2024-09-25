FROM golang:1.22.6-alpine3.20

WORKDIR /app

RUN go install github.com/air-verse/air@v1.52.3

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

CMD ["air"]
