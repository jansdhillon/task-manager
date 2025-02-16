# syntax=docker/dockerfile:1
FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/task-manager/*.go ./

EXPOSE 8080

RUN --mount=type=secret,id=github_token,env=GITHUB_TOKEN

RUN --mount=type=secret,id=supabase_anon_key,env=SUPABASE_ANON_KEY

RUN --mount=type=secret,id=supabase_service_role_key,env=SUPABASE_SERVICE_ROLE_KEY

RUN --mount=type=secret,id=supabase_url,env=SUPABASE_URL

RUN CGO_ENABLED=0 GOOS=linux go build -o task-manager main.go

CMD ["./task-manager"]
