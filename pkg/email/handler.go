package email

import (
	"context"
	"errors"
	"os"
	"time"

	pb "github.com/Helvind/email-service/proto"
	sp "github.com/SparkPost/gosparkpost"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendGridProvider impl of provider
type SendGridProvider struct {
	client *sendgrid.Client
}

func (p *SendGridProvider) send(email *pb.Email) (*pb.SendEmailResponse, error) {
	message := mail.NewSingleEmail(
		&mail.Email{Address: email.From},
		email.Subject,
		&mail.Email{Address: email.To},
		email.Body,
		email.Body,
	)
	response, err := p.client.Send(message)
	return &pb.SendEmailResponse{
		Status: int32(response.StatusCode),
	}, err
}

func (p *SendGridProvider) init() {
	p.client = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
}

// MailGunProvider impl of provider
type MailGunProvider struct {
	client *mailgun.MailgunImpl
}

func (p *MailGunProvider) send(email *pb.Email) (*pb.SendEmailResponse, error) {
	m := p.client.NewMessage(
		email.From,
		email.Subject,
		email.Body,
		email.To,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := p.client.Send(ctx, m)
	if err != nil {
		return &pb.SendEmailResponse{
			Status: 500,
		}, err
	}
	return &pb.SendEmailResponse{
		Status: 202,
	}, nil
}

func (p *MailGunProvider) init() {
	p.client = mailgun.NewMailgun(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
	)
}

// SparkPostProvider impl of provider
type SparkPostProvider struct {
	client *sp.Client
}

func (p *SparkPostProvider) send(email *pb.Email) (*pb.SendEmailResponse, error) {
	tx := &sp.Transmission{
		Recipients: []string{email.To},
		Content: sp.Content{
			HTML:    email.Body,
			From:    email.From,
			Subject: email.Subject,
		},
	}
	_, res, err := p.client.Send(tx)
	return &pb.SendEmailResponse{
		Status: int32(res.HTTP.StatusCode),
	}, err
}

func (p *SparkPostProvider) init() {
	apiKey := os.Getenv("SPARKPOST_API_KEY")
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     apiKey,
		ApiVersion: 1,
	}
	var client sp.Client
	client.Init(cfg)
}

// FakeProvider for testing purposes only
type FakeProvider struct{}

func (p *FakeProvider) send(email *pb.Email) (*pb.SendEmailResponse, error) {
	time.Sleep(time.Millisecond * 100)
	return &pb.SendEmailResponse{
		Status: 202,
	}, nil
}

func (p *FakeProvider) init() {}

// ValidateEmail validates the fields of the Email
func ValidateEmail(email *pb.Email) error {
	if email.GetBody() == "" {
		return errors.New("missing body")
	}
	if email.GetFrom() == "" {
		return errors.New("missing from field")
	}
	if email.GetTo() == "" {
		return errors.New("missing to field")
	}
	if email.GetSubject() == "" {
		return errors.New("missing subject")
	}
	return nil
}

// SendEmail sends the provided Email through the active provider set on the Service
func (e *Service) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (
	*pb.SendEmailResponse, error) {
	if err := ValidateEmail(req.GetEmail()); err != nil {
		return &pb.SendEmailResponse{
			Status: 400,
		}, status.Errorf(codes.InvalidArgument, "Invalid E-mail, %s", err)
	}

	provider := e.providers[e.active]
	if provider == nil {
		zap.L().Warn("No provider match active", zap.String("active provider", e.active))
		return nil, status.Errorf(codes.Internal, "Could not find provider")
	}

	zap.L().Info("Sending email", zap.String("provider", e.active))

	resp, err := provider.send(req.GetEmail())
	if err != nil {
		zap.L().Warn("Could not send E-mail", zap.Error(err), zap.String("active provider", e.active))
		return &pb.SendEmailResponse{
			Status: 500,
		}, status.Errorf(codes.Unavailable, "Could not send E-mail")
	}

	return resp, nil
}
