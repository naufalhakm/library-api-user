migration_up: 
	migrate -path ./pkg/database/migration/ -database "postgresql://user:password@localhost:5432/sweatsparks?sslmode=disable" -verbose up

migration_down: 
	migrate -path ./pkg/database/migration/ -database "postgresql://user:password@localhost:5432/sweatsparks?sslmode=disable" -verbose down

gen_proto:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/book/book.proto
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/auth/auth.proto
