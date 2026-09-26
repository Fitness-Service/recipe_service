module github.com/Fitness-Service/recipe-service

go 1.25.0

require (
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.12.3
)

require (
	github.com/Fitness-Service/shared v0.0.0-20260629195911-aa9f7bcaad3f // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/grpc v1.84.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/Fitness-Service/shared => ../shared
