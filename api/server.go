package api

import (
	"context"
	"fmt"
	"github.com/adam-bunce/grpc-todo/controllers"
	todo_service "github.com/adam-bunce/grpc-todo/domain/proto"
	logr "github.com/adam-bunce/grpc-todo/helpers"
	"github.com/adam-bunce/grpc-todo/repository"
	"github.com/adam-bunce/grpc-todo/usecase"
	"github.com/adam-bunce/grpc-todo/variables"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func CreateServer() {
	port := variables.GlobalConfig.ServerPort

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor))

	// register stuff to server
	todoRepo := repository.NewTaskRepository(variables.DB)
	todoUseCase := usecase.NewTodoUseCase(todoRepo)

	todoServiceServer := controllers.TodoServiceServer{UseCase: todoUseCase}

	todo_service.RegisterTodoServiceServer(server, &todoServiceServer)
	reflection.Register(server)

	// start grpc on a diff "thread"
	go func() {
		logr.Info(fmt.Sprintf("Server running on :%d", port))
		if err := server.Serve(lis); err != nil {
			logr.Error(fmt.Sprintf("Error starting server: %v", err))
			panic(err)
		}
	}()

	// below is setup for http part
	// dial in to grpc services that queries are gonna be proxied to
	conn, err := grpc.DialContext(context.TODO(), "127.0.0.1:8080", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logr.Error(fmt.Sprintf("Error dialing server: %v", err))
		panic(err)
	}

	gwmux := runtime.NewServeMux()

	// add cors middleware
	// this is the first thing that hits
	// handler wrapping handler..
	// https://gist.github.com/nedimf/47f1a4f295f46601547fde55e48203aa
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// catch options request
		// this intercepts option's requests for preflight
		// preflight is when a browser has a funky request and needs to make sure it's allowed before sending
		// sinple get request's dont' do preflight
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "content-type")
			w.WriteHeader(204)
			// finish
			return
		}

		if origin := r.Header.Get("Origin"); origin != "" {
			fmt.Println("request from: ", origin)
			//origins := "http:ligma.com,http://localhost:8000" // can't actually have two
			origins := "*" // allow all
			// could also use origin which is the origin of the requesting server / ~browser
			// browsers auto add origin headers

			w.Header().Set("Access-Control-Allow-Origin", origins)
		}
		// this calls the next part of the handler chain
		gwmux.ServeHTTP(w, r)
	})

	err = todo_service.RegisterTodoServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		logr.Error(fmt.Sprintf("failed to register service handler: %v", err))
		panic(err)
	}

	gwServer := &http.Server{
		Addr:    ":9000",
		Handler: handler,
	}
	logr.Info("Serving gRPC-Gateway on http://0.0.0.0:9000") // 1337 is actually busy wtf
	log.Fatalln(gwServer.ListenAndServe())
}
