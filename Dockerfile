FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/task-manager/*.go ./

EXPOSE 8080

RUN CGO_ENABLED=0 GOOS=linux go build -o task-manager main.go

CMD ["./task-manager"]
