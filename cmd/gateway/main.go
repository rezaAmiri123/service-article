package main

import (
	"context"
	"github.com/rezaAmiri123/service-article/pkg/utils"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rezaAmiri123/service-user/cmd/config"
	pb "github.com/rezaAmiri123/service-article/gen/pb"
	"google.golang.org/grpc"
)

func run(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ropts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	}

	mux := runtime.NewServeMux(ropts...)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterArticlesHandlerFromEndpoint(context.Background(), mux, cfg.Gateway.GetServerAddress(), opts)
	if err != nil {
		return err
	}
	log.Printf("starting gateway server on port %v", cfg.Gateway.Port)
	return http.ListenAndServe(cfg.Gateway.Port, mux)
}

func main() {
	log.Println("Starting user gateway microservice")

	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	if err := run(cfg); err != nil {
		logrus.Fatal(err.Error())
	}
}
