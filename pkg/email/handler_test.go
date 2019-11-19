package email

import (
	"testing"

	pb "github.com/Helvind/email-service/proto"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    *pb.Email
		expected bool
	}{
		{&pb.Email{To: "to", From: "from", Subject: "Sub", Body: "Body"}, true},
		{&pb.Email{To: "", From: "from", Subject: "Sub", Body: "Body"}, false},
		{&pb.Email{To: "to", From: "", Subject: "Sub", Body: "Body"}, false},
		{&pb.Email{To: "to", From: "from", Subject: "", Body: "Body"}, false},
		{&pb.Email{To: "to", From: "from", Subject: "Sub", Body: ""}, false},
	}
	for _, e := range tests {
		err := ValidateEmail(e.email)
		if err != nil && e.expected || err == nil && !e.expected {
			t.Fail()
		}
	}
}
