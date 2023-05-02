to run do elm reactor in root

then run build cmd/main.go, everything should be there, no env vars



generate basic grpc proto file
```bash
protoc -I .\domain\proto\ --go_out .\domain\proto\ --go_opt paths=source_relative --go-grpc_out .\domain\proto\ --go-grpc_opt paths=source_relative .\domain\proto\todo_service.pr
oto
```

```bash
protoc -I .\domain\proto\ --go_out .\domain\proto\ --go_opt paths=source_relative --go-grpc_out .\domain\proto\ --go-grpc_opt paths=source_relative .\domain\proto\todo_service.proto --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative .\domain\proto\todo_service.proto
```

```
protoc -I .\domain\proto\ --go_out .\domain\proto\ --go_opt paths=source_relative --go-grpc_out .\domain\proto\ --go-grpc_opt paths=source_relative --grpc-gateway_out .\domain\pr
oto\ --grpc-gateway_opt paths=source_relative .\domain\proto\todo_service.proto
```
