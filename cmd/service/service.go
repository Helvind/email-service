package main

import (
	"log"
	"net"
	"net/http"

	"github.com/Helvind/email-service/pkg/email"
	pb "github.com/Helvind/email-service/proto"
	health "github.com/Helvind/email-service/pkg/health/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	api_health "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	port = ":50051"
)

// starts the Prometheus stats endpoint server
func startPromHTTPServer(promport string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":"+promport, nil); err != nil {
		log.Println("prometheus err", promport)
	}
}

func main() {
	// Init logging
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger, %q", err.Error())
	}
	zap.ReplaceGlobals(logger)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		zap.L().Fatal("failed to listen", zap.Error(err))
	}

	go startPromHTTPServer("5001")

	s := grpc.NewServer()

	// Register: Health
	healthServ := health.NewHealthCheckService()
	api_health.RegisterHealthServer(s, healthServ)

	// List of available providers
	// TODO(helvind) Load this form config
	providers := map[string]email.Provider{
		"sendgrid":  &email.SendGridProvider{},
		"mailgun":   &email.MailGunProvider{},
		"sparkpost": &email.SparkPostProvider{},
		"fake":      &email.FakeProvider{},
	}

	pb.RegisterEmailSenderServer(s, email.NewService("fake", providers))
	if err := s.Serve(lis); err != nil {
		zap.L().Fatal("failed to serve: %v", zap.Error(err))
	}
}
