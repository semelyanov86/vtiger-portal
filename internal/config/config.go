package config

import (
	_ "github.com/octoper/go-ray"
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
		Email EmailConfig `yaml:"email"`
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
