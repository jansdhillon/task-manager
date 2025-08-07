FROM golang:bookworm

RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN ./generate.sh

RUN CGO_ENABLED=0 GOOS=linux go build -o /task-manager ./cmd/task-manager

EXPOSE 8080

CMD ["/task-manager"]
