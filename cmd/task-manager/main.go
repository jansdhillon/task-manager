package main

import (
	"fmt"
	"net/http"
	"os"
)

type SupabaseClient struct {
	SupabaseUrl            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string
}

func getTasksHandler(client *SupabaseClient) http.HandlerFunc {
	_, err := http.Get(client.SupabaseUrl)
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
}

func rootHandler(client *SupabaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		fmt.Fprintln(w, "Task Manager")
	}

}

func main() {

	client := &SupabaseClient{
		SupabaseUrl:            getEnv("SUPABASE_URL"),
		SupabaseAnonKey:        getEnv("SUPABASE_ANON_KEY"),
		SupabaseServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY"),
	}

	http.HandleFunc("/", rootHandler(client))
	http.HandleFunc("/tasks", getTasksHandler((client)))

	fmt.Printf("Server starting on port 8080\n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Env var %s not found, exiting...", key)
		os.Exit(1)
	}
	return value
}
