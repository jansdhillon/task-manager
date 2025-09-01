FROM golang:bookworm
ARG GCP_PROJECT_ID

ENV GCP_PROJECT_ID=${GCP_PROJECT_ID}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /task-manager ./cmd/task-manager

EXPOSE 8080

CMD ["/task-manager"]
