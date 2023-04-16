package config

import (
	_ "github.com/octoper/go-ray"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type (
	Config struct {
		Environment string         `yaml:"environment"`
		HTTP        HTTPConfig     `yaml:"http"`
		Cache       CacheConfig    `yaml:"cache"`
		Db          DatabaseConfig `yaml:"db"`
		Smtp        SmtpConfig     `yaml:"smtp"`
		Limiter     Limiter        `yaml:"limiter"`
		Cors        struct {
			TrustedOrigins []string `yaml:"trustedOrigins"`
		}
		Email  EmailConfig  `yaml:"email"`
		Vtiger VtigerConfig `yaml:"vtiger"`
	}
	HTTPConfig struct {
		Host               string        `yaml:"host"`
		Port               int           `yaml:"port"`
		ReadTimeout        time.Duration `yaml:"readTimeout"`
		WriteTimeout       time.Duration `yaml:"writeTimeout"`
		MaxHeaderMegabytes int           `yaml:"maxHeaderBytes"`
	}
	CacheConfig struct {
		TTL time.Duration `yaml:"ttl"`
	}
	DatabaseConfig struct {
		Dsn          string
		Host         string
		Login        string
		Password     string
		Dbname       string
		MaxOpenConns int    `yaml:"maxOpenConns"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxIdleTime  string `yaml:"maxIdleTime"`
	}
	SmtpConfig struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
		TTL     time.Duration
	}
	EmailConfig struct {
		Templates EmailTemplates `yaml:"templates"`
		Subjects  EmailSubjects  `yaml:"subjects"`
	}

	EmailTemplates struct {
		RegistrationEmail string `yaml:"registrationEmail"`
		TicketSuccessful  string `yaml:"ticketSuccessful"`
	}

	EmailSubjects struct {
		RegistrationEmail string `yaml:"registrationEmail"`
		TicketSuccessful  string `yaml:"ticketSuccessful"`
	}
	VtigerConfig struct {
		Connection vtiger.VtigerConnectionConfig `yaml:"connection"`
		Business   VtigerBusinessConfig          `yaml:"business"`
	}
	VtigerBusinessConfig struct {
		EmailField string `yaml:"emailField"`
		CodeField  string `yaml:"codeField"`
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) *Config {
	var cfg *Config
	cfg = readConfigFile(cfg, configsDir)
	cfg.Db.Dsn = cfg.Db.Login + ":" + cfg.Db.Password + "@" + cfg.Db.Host + "/" + cfg.Db.Dbname + "?parseTime=true"
	return cfg
}

func readConfigFile(cfg *Config, configsDir string) *Config {
	bytesOut, err := os.ReadFile(configsDir + "/portal.yaml")

	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bytesOut, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
