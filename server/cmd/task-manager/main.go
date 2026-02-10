package main

import (
	"fmt"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	taskv1connect "github.com/jansdhillon/task-manager/proto/gen/task/v1/taskv1connect"
	"github.com/jansdhillon/task-manager/server/internal/db"
	"github.com/jansdhillon/task-manager/server/internal/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	taskDB := db.NewTaskDB()
	svc := server.NewTaskServer(taskDB)
	mux := http.NewServeMux()
	path, handler := taskv1connect.NewTaskServiceHandler(
		svc,
		connect.WithInterceptors(validate.NewInterceptor()),
	)
	mux.Handle(path, handler)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)
	port := os.Getenv("TASK_MANAGER_PORT")
	if port == "" {
		port = "8080"
	}
	address := os.Getenv("TASK_MANAGER_ADDRESS")
	if address == "" {
		address = "0.0.0.0"
	}
	s := http.Server{
		Addr:      fmt.Sprintf("%s:%s", address, port),
		Handler:   mux,
		Protocols: p,
	}
	s.ListenAndServe()

}
