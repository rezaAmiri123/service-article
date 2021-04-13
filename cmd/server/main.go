package main

import (
	"github.com/rezaAmiri123/service-article/cmd/config"
	"github.com/rezaAmiri123/service-article/internal/handler"
	"github.com/rezaAmiri123/service-article/internal/model"
	"github.com/rezaAmiri123/service-article/internal/repository"
	"github.com/rezaAmiri123/service-article/pkg/jaeger"
	"github.com/rezaAmiri123/service-article/pkg/logger"
	"github.com/rezaAmiri123/service-article/pkg/mysql"
	"github.com/rezaAmiri123/service-article/pkg/utils"
	userPb "github.com/rezaAmiri123/service-user/gen/pb"

	"log"
	"net"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "github.com/rezaAmiri123/service-article/gen/pb"
)

func main() {
	log.Println("Starting article server microservice")

	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Server.AppVersion)

	db := mysql.NewGormDB(cfg)
	defer db.Close()
	model.AutoMigrate(db)


	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	repo := repository.NewORMArticleRepository(db)

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(cfg.UserServer.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	userConn := userPb.NewUsersClient(conn)


	h := handler.NewArticleHandler(repo, appLogger,userConn)
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	srv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  cfg.Server.MaxConnectionAge * time.Minute,
		Time:              cfg.Server.Timeout * time.Minute,
	}),
		//grpc.UnaryInterceptor(im.Logger),
		//grpc.ChainUnaryInterceptor(
		//	grpc_ctxtags.UnaryServerInterceptor(),
		//	grpc_prometheus.UnaryServerInterceptor,
		//	grpcrecovery.UnaryServerInterceptor(),
		//),
	)

	pb.RegisterArticlesServer(srv, h)
	appLogger.Info("server starts at ", cfg.Server.Port)

	if err := srv.Serve(lis); err != nil {
		appLogger.Fatal(err.Error())
	}

}
