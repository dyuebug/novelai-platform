package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Auth      AuthConfig
	OAuth     OAuthConfig
	Email     EmailConfig
	RateLimit RateLimitConfig
	GRPC      GRPCConfig
	Log       LogConfig
}

type ServerConfig struct {
	Addr string
	Mode string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
	LogLevel     string
}

type AuthConfig struct {
	JWTSecret string
	Issuer    string
	Audience  string
}

// OAuthConfig OAuth 配置
type OAuthConfig struct {
	BaseURL  string // 应用基础 URL，用于回调
	GitHub   OAuthProviderConfig
	Google   OAuthProviderConfig
	LinuxDo  OAuthProviderConfig
}

// OAuthProviderConfig OAuth 提供商配置
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	Enabled      bool
}

type RateLimitConfig struct {
	RPS   float64
	Burst int
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromAddress  string
	FromName     string
	Enabled      bool
}

type GRPCConfig struct {
	AIServiceAddr string
	Timeout       time.Duration
	MaxConns      int
	RetryMax      int
	RetryBackoff  time.Duration
	// TLS 配置
	TLSEnabled bool
	CertFile   string // 客户端证书
	KeyFile    string // 客户端私钥
	CAFile     string // CA 证书
}

type LogConfig struct {
	Level string
}

func Load() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	setDefaults(v)
	_ = v.ReadInConfig()

	return &Config{
		Server: ServerConfig{
			Addr: v.GetString("server.addr"),
			Mode: v.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Host:         v.GetString("database.host"),
			Port:         v.GetInt("database.port"),
			User:         v.GetString("database.user"),
			Password:     v.GetString("database.password"),
			DBName:       v.GetString("database.dbname"),
			SSLMode:      v.GetString("database.sslmode"),
			MaxIdleConns: v.GetInt("database.max_idle_conns"),
			MaxOpenConns: v.GetInt("database.max_open_conns"),
			LogLevel:     v.GetString("database.log_level"),
		},
		Auth: AuthConfig{
			JWTSecret: v.GetString("auth.jwt_secret"),
			Issuer:    v.GetString("auth.issuer"),
			Audience:  v.GetString("auth.audience"),
		},
		OAuth: OAuthConfig{
			BaseURL: v.GetString("oauth.base_url"),
			GitHub: OAuthProviderConfig{
				ClientID:     v.GetString("oauth.github.client_id"),
				ClientSecret: v.GetString("oauth.github.client_secret"),
				RedirectURL:  v.GetString("oauth.github.redirect_url"),
				Scopes:       v.GetStringSlice("oauth.github.scopes"),
				AuthURL:      "https://github.com/login/oauth/authorize",
				TokenURL:     "https://github.com/login/oauth/access_token",
				UserInfoURL:  "https://api.github.com/user",
				Enabled:      v.GetBool("oauth.github.enabled"),
			},
			Google: OAuthProviderConfig{
				ClientID:     v.GetString("oauth.google.client_id"),
				ClientSecret: v.GetString("oauth.google.client_secret"),
				RedirectURL:  v.GetString("oauth.google.redirect_url"),
				Scopes:       v.GetStringSlice("oauth.google.scopes"),
				AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL:     "https://oauth2.googleapis.com/token",
				UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
				Enabled:      v.GetBool("oauth.google.enabled"),
			},
			LinuxDo: OAuthProviderConfig{
				ClientID:     v.GetString("oauth.linuxdo.client_id"),
				ClientSecret: v.GetString("oauth.linuxdo.client_secret"),
				RedirectURL:  v.GetString("oauth.linuxdo.redirect_url"),
				Scopes:       v.GetStringSlice("oauth.linuxdo.scopes"),
				AuthURL:      v.GetString("oauth.linuxdo.auth_url"),
				TokenURL:     v.GetString("oauth.linuxdo.token_url"),
				UserInfoURL:  v.GetString("oauth.linuxdo.userinfo_url"),
				Enabled:      v.GetBool("oauth.linuxdo.enabled"),
			},
		},
		Email: EmailConfig{
			SMTPHost:     v.GetString("email.smtp_host"),
			SMTPPort:     v.GetInt("email.smtp_port"),
			SMTPUsername: v.GetString("email.smtp_username"),
			SMTPPassword: v.GetString("email.smtp_password"),
			FromAddress:  v.GetString("email.from_address"),
			FromName:     v.GetString("email.from_name"),
			Enabled:      v.GetBool("email.enabled"),
		},
		RateLimit: RateLimitConfig{
			RPS:   v.GetFloat64("rate_limit.rps"),
			Burst: v.GetInt("rate_limit.burst"),
		},
		GRPC: GRPCConfig{
			AIServiceAddr: v.GetString("grpc.ai_service_addr"),
			Timeout:       time.Duration(v.GetInt("grpc.timeout_ms")) * time.Millisecond,
			MaxConns:      v.GetInt("grpc.max_conns"),
			RetryMax:      v.GetInt("grpc.retry_max"),
			RetryBackoff:  time.Duration(v.GetInt("grpc.retry_backoff_ms")) * time.Millisecond,
			TLSEnabled:    v.GetBool("grpc.tls_enabled"),
			CertFile:      v.GetString("grpc.cert_file"),
			KeyFile:       v.GetString("grpc.key_file"),
			CAFile:        v.GetString("grpc.ca_file"),
		},
		Log: LogConfig{
			Level: v.GetString("log.level"),
		},
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("server.mode", "release")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.dbname", "novelai")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.log_level", "silent")
	v.SetDefault("auth.jwt_secret", "change-me")
	v.SetDefault("auth.issuer", "novelai")
	v.SetDefault("auth.audience", "novelai-client")
	v.SetDefault("oauth.base_url", "http://localhost:8080")
	v.SetDefault("oauth.github.enabled", false)
	v.SetDefault("oauth.github.scopes", []string{"user:email"})
	v.SetDefault("oauth.google.enabled", false)
	v.SetDefault("oauth.google.scopes", []string{"openid", "email", "profile"})
	v.SetDefault("oauth.linuxdo.enabled", false)
	v.SetDefault("oauth.linuxdo.auth_url", "https://connect.linux.do/oauth2/authorize")
	v.SetDefault("oauth.linuxdo.token_url", "https://connect.linux.do/oauth2/token")
	v.SetDefault("oauth.linuxdo.userinfo_url", "https://connect.linux.do/api/user")
	v.SetDefault("oauth.linuxdo.scopes", []string{"user:email"})
	v.SetDefault("email.enabled", false)
	v.SetDefault("email.smtp_host", "smtp.gmail.com")
	v.SetDefault("email.smtp_port", 587)
	v.SetDefault("email.from_name", "NovelAI")
	v.SetDefault("rate_limit.rps", 20)
	v.SetDefault("rate_limit.burst", 40)
	v.SetDefault("grpc.ai_service_addr", "127.0.0.1:50051")
	v.SetDefault("grpc.timeout_ms", 30000)
	v.SetDefault("grpc.max_conns", 8)
	v.SetDefault("grpc.retry_max", 2)
	v.SetDefault("grpc.retry_backoff_ms", 200)
	v.SetDefault("grpc.tls_enabled", false)
	v.SetDefault("grpc.cert_file", "")
	v.SetDefault("grpc.key_file", "")
	v.SetDefault("grpc.ca_file", "")
	v.SetDefault("log.level", "info")
}
