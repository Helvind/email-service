package email

import (
	pb "github.com/Helvind/email-service/proto"
)

// Provider is a interface for different e-mail sending providers
type Provider interface {
	init()
	send(*pb.Email) (*pb.SendEmailResponse, error)
}

// Service can send Emails to different providers
type Service struct {
	active    string
	providers map[string]Provider
}

// NewService returns a new EmailService
func NewService(active string, providers map[string]Provider) *Service {
	for _, p := range providers {
		p.init()
	}

	return &Service{
		active:    active,
		providers: providers,
	}
}
