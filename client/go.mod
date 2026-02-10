module github.com/jansdhillon/task-manager/client

go 1.24.0

require (
	connectrpc.com/connect v1.19.1
	github.com/jansdhillon/task-manager/proto v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v3 v3.3.8
	google.golang.org/grpc v1.78.0
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/jansdhillon/task-manager/proto => ../proto
