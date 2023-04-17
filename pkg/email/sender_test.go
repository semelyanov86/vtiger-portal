package email

import (
	"strings"
	"testing"
)

func TestSendEmailInput_GenerateBodyFromHTML(t *testing.T) {
	input := SendEmailInput{}
	templateFileName := "smtp/templates/registration_email.html"
	data := struct {
		Name string
	}{
		Name: "John",
	}

	err := input.GenerateBodyFromHTML(templateFileName, data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedBody := "Thank You for Registering\n"
	if strings.Contains(input.Body, expectedBody) {
		t.Errorf("Expected body to be %q, but got %q", expectedBody, input.Body)
	}
}
