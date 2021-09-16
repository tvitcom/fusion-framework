package config

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-env"
	"github.com/tvitcom/fusion-framework/pkg/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"errors"
	"flag"
)

const (
	defaultJWTExpirationHours = 2
)
var configFile = flag.String("config", "./configs/dev.yml", "path to the config file")

// Config represents an application configuration.
type Config struct {
	AppMode string      `yaml:"app_mode" env:"APP_MODE"`
	AppName  string      `yaml:"app_name" env:"APP_NAME"`
	DSN string           `yaml:"dsn" env:"DSN"`
	DBType string        `yaml:"db_type" env:"DB_TYPE"`
	JWTSigningKey string `yaml:"jwt_signing_key" env:"JWT_SIGNING_KEY"`
	JWTExpiration int    `yaml:"jwt_expiration" env:"JWT_EXPIRATION"`
	MailSmtphost string  `yaml:"mail_smtphost" env:"MAIL_SMTPHOST"`
	MailSmtpport string  `yaml:"mail_smtpport" env:"MAIL_SMTPPORT"`
	MailUsername string  `yaml:"mail_username" env:"MAIL_USERNAME"`
	MailPassword string  `yaml:"mail_password" env:"MAIL_PASSWORD"`
	AppFqdn string       `yaml:"app_fqdn" env:"APP_FQDN"`
	HttpEntrypoint string       `yaml:"http_entrypoint" env:"HTTP_ENTRYPOINt"`
	WebservName string          `yaml:"webserv_name" env:"WEBSERV_NAME"`
	GoogleCredentialFile string `yaml:"google_credential_file" env:"GOOGLE_CREDENTIAL_FILE"`
	GoogleRedirectPath string   `yaml:"google_redirect_path" env:"GOOGLE_REDIRECT_PATH"`
	AppSecretKey string  `yaml:"app_secret_key" env:"APP_SECRET_KEY"`
	BizName string       `yaml:"biz_name" env:"BIZ_NAME"`
	BizShortname string  `yaml:"biz_shortname" env:"BIZ_SHORTNAME"`
	BizEmail string      `yaml:"biz_email" env:"BIZ_EMAIL"`
	BizPhone string      `yaml:"biz_phone" env:"BIZ_PHONE"`
	BizPhone2 string     `yaml:"biz_phone2" env:"BIZ_PHONE2"`
	BizLogo string       `yaml:"biz_logo" env:"BIZ_LOGO"`
}

func init() {
	flag.Parse()
}

// Validate validates the application configuration.
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AppMode, validation.Required),
		validation.Field(&c.AppName, validation.Required),
		validation.Field(&c.DBType, validation.Required),
		validation.Field(&c.DSN, validation.Required),
		validation.Field(&c.JWTSigningKey, validation.Required),
		validation.Field(&c.JWTExpiration, validation.Required),
		validation.Field(&c.MailSmtphost, validation.Required),
		validation.Field(&c.MailSmtpport, validation.Required),
		validation.Field(&c.MailUsername, validation.Required),
		validation.Field(&c.MailPassword, validation.Required),
		validation.Field(&c.AppFqdn, validation.Required),
		validation.Field(&c.HttpEntrypoint, validation.Required),
		validation.Field(&c.WebservName, validation.Required),
		validation.Field(&c.GoogleCredentialFile, validation.Required),
		validation.Field(&c.GoogleRedirectPath, validation.Required),
		validation.Field(&c.AppSecretKey, validation.Required),
		validation.Field(&c.BizName, validation.Required),
		validation.Field(&c.BizShortname, validation.Required),
		validation.Field(&c.BizEmail, validation.Required),
		validation.Field(&c.BizPhone),
		validation.Field(&c.BizPhone2),
		validation.Field(&c.BizLogo, validation.Required),
	)
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(logger log.Logger) (*Config, error) {
	// default config
	c := Config{
		JWTExpiration: defaultJWTExpirationHours,
	}
	// load from YAML config file
	bytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, errors.New("config file fs read failed")
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	// load from environment variables prefixed with "APP_"
	err = env.New("APP_", logger.Infof).Load(&c)
	if err != nil {
		return nil, err
	}

	// validation
	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}

func GetPort(ep string) (string, error) {
	if strings.Index(ep, ":") > 0 {
		return strings.SplitN(ep, ":",2)[1], nil
	}
	return "", errors.New("Absent port number in Entry point string")
}
