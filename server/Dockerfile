FROM golang:bookworm
ARG GCP_PROJECT_ID
ARG POSTGRES_HOST
ARG POSTGRES_USER

ENV GCP_PROJECT_ID=${GCP_PROJECT_ID}
ENV POSTGRES_HOST=${POSTGRES_HOST}
ENV POSTGRES_USER=${POSTGRES_USER}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /task-manager ./cmd/task-manager

EXPOSE 8080

CMD ["/task-manager"]
