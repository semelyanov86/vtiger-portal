package mock

import "github.com/semelyanov86/vtiger-portal/pkg/email"
import "github.com/stretchr/testify/mock"

type EmailProvider struct {
	mock.Mock
}

func (m *EmailProvider) AddEmailToList(inp email.AddEmailInput) error {
	args := m.Called(inp)

	return args.Error(0)
}

type EmailSender struct {
	mock.Mock
}

func (m *EmailSender) Send(recipient, templateFile string, data any) error {
	args := m.Called(recipient)

	return args.Error(0)
}
