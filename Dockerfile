FROM golang:1.24-bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /task-manager ./cmd/task-manager

EXPOSE 8080

CMD ["/task-manager"]
