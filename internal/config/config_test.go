package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Set up test configuration file
	testDir := "testdata"
	testFile := "portal.yaml"
	testContent := `
environment: test
http:
  host: localhost
  port: 8080
  readTimeout: 10s
  writeTimeout: 10s
  maxHeaderBytes: 1024
cache:
  ttl: 1h
db:
  host: localhost
  login: user
  password: password
  dbname: testdb
  maxOpenConns: 10
  maxIdleConns: 5
  maxIdleTime: 10m
smtp:
  host: smtp.gmail.com
  port: 587
  username: user
  password: password
  sender: sender@example.com
limiter:
  rps: 10
  burst: 100
  enabled: true
cors:
  trustedOrigins:
    - http://localhost:3000
email:
  templates:
    registrationEmail: test-template.html
    ticketSuccessful: test-template.html
  subjects:
    registrationEmail: Test Registration Email
    ticketSuccessful: Test Ticket Successful
vtiger:
  connection:
    url: "https://serv.itvolga.com"
    login: "admin"
    password: ""
`
	err := os.Mkdir(testDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(testDir+"/"+testFile, []byte(testContent), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	// Call Init function and check if configuration is populated correctly
	cfg := Init(testDir)
	assert.Equal(t, "test", cfg.Environment)
	assert.Equal(t, "localhost", cfg.HTTP.Host)
	assert.Equal(t, 8080, cfg.HTTP.Port)
	assert.Equal(t, 10*time.Second, cfg.HTTP.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.HTTP.WriteTimeout)
	assert.Equal(t, 1024, cfg.HTTP.MaxHeaderMegabytes)
	assert.Equal(t, 1*time.Hour, cfg.Cache.TTL)
	assert.Equal(t, "user:password@localhost/testdb?parseTime=true", cfg.Db.Dsn)
	assert.Equal(t, 10, cfg.Db.MaxOpenConns)
	assert.Equal(t, 5, cfg.Db.MaxIdleConns)
	assert.Equal(t, "10m", cfg.Db.MaxIdleTime)
	assert.Equal(t, "smtp.gmail.com", cfg.Smtp.Host)
	assert.Equal(t, 587, cfg.Smtp.Port)
	assert.Equal(t, "user", cfg.Smtp.Username)
	assert.Equal(t, "password", cfg.Smtp.Password)
	assert.Equal(t, "sender@example.com", cfg.Smtp.Sender)
	assert.Equal(t, 10.0, cfg.Limiter.Rps)
	assert.Equal(t, 100, cfg.Limiter.Burst)
	assert.Equal(t, true, cfg.Limiter.Enabled)
	assert.Equal(t, []string{"http://localhost:3000"}, cfg.Cors.TrustedOrigins)
	assert.Equal(t, "test-template.html", cfg.Email.Templates.RegistrationEmail)
	assert.Equal(t, "test-template.html", cfg.Email.Templates.TicketSuccessful)
	assert.Equal(t, "Test Registration Email", cfg.Email.Subjects.RegistrationEmail)
	assert.Equal(t, "Test Ticket Successful", cfg.Email.Subjects.TicketSuccessful)
}
