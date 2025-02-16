FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/task-manager/*.go ./

EXPOSE 8080

RUN --mount=type=secret,id=github_token,env=GITHUB_TOKEN

RUN --mount=type=secret,id=github_token,env=SUPABASE_URL

RUN --mount=type=secret,id=github_token,env=SUPABASE_ANON_KEY

RUN --mount=type=secret,id=github_token,env=SUPABASE_SERVICE_ROLE_KEY

RUN CGO_ENABLED=0 GOOS=linux go build -o task-manager main.go

CMD ["./task-manager"]
