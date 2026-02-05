package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
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

	// 内部字段：记录配置来源
	configSources map[string]string
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

	cfg := &Config{
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

	// 环境变量优先级：覆盖 config.yaml 的配置
	applyEnvOverrides(cfg)

	// 验证必需配置项
	validateConfig(cfg)

	// 输出配置来源日志
	logConfigSources(cfg)

	return cfg
}

// parsePostgresURL 解析 PostgreSQL 连接字符串
// 格式: postgres://user:password@host:port/dbname?sslmode=disable
func parsePostgresURL(dbURL string) (*DatabaseConfig, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_URL format: %w", err)
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return nil, fmt.Errorf("invalid DATABASE_URL scheme: expected postgres or postgresql, got %s", u.Scheme)
	}

	cfg := &DatabaseConfig{
		Host:   u.Hostname(),
		DBName: strings.TrimPrefix(u.Path, "/"),
	}

	// 解析端口
	if u.Port() != "" {
		var port int
		fmt.Sscanf(u.Port(), "%d", &port)
		cfg.Port = port
	} else {
		cfg.Port = 5432 // 默认 PostgreSQL 端口
	}

	// 解析用户名和密码
	if u.User != nil {
		cfg.User = u.User.Username()
		if password, ok := u.User.Password(); ok {
			cfg.Password = password
		}
	}

	// 解析查询参数
	query := u.Query()
	if sslmode := query.Get("sslmode"); sslmode != "" {
		cfg.SSLMode = sslmode
	} else {
		cfg.SSLMode = "disable"
	}

	return cfg, nil
}

// applyEnvOverrides 应用环境变量覆盖
func applyEnvOverrides(cfg *Config) {
	configSources := make(map[string]string)

	// DATABASE_URL 优先级最高
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dbCfg, err := parsePostgresURL(dbURL)
		if err != nil {
			log.Printf("WARNING: Failed to parse DATABASE_URL: %v, using config.yaml values", err)
		} else {
			cfg.Database.Host = dbCfg.Host
			cfg.Database.Port = dbCfg.Port
			cfg.Database.User = dbCfg.User
			cfg.Database.Password = dbCfg.Password
			cfg.Database.DBName = dbCfg.DBName
			cfg.Database.SSLMode = dbCfg.SSLMode
			configSources["database"] = "DATABASE_URL"
		}
	} else {
		configSources["database"] = "config.yaml"
	}

	// GRPC_AI_SERVICE_ADDR 环境变量
	if grpcAddr := os.Getenv("GRPC_AI_SERVICE_ADDR"); grpcAddr != "" {
		cfg.GRPC.AIServiceAddr = grpcAddr
		configSources["grpc"] = "GRPC_AI_SERVICE_ADDR"
	} else {
		configSources["grpc"] = "config.yaml"
	}

	// JWT_SECRET 环境变量
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.Auth.JWTSecret = jwtSecret
		configSources["jwt"] = "JWT_SECRET"

		// 验证 JWT_SECRET 长度
		if len(jwtSecret) < 32 {
			log.Printf("WARNING: JWT_SECRET is shorter than recommended 32 characters (current: %d). Consider using a stronger secret.", len(jwtSecret))
		}
	} else {
		configSources["jwt"] = "config.yaml"
	}

	// 存储配置来源供日志使用
	cfg.configSources = configSources
}

// validateConfig 验证必需配置项
func validateConfig(cfg *Config) {
	var errors []string

	// 验证数据库配置
	if cfg.Database.Host == "" {
		errors = append(errors, "Database host is not configured (set DATABASE_URL or database.host in config.yaml)")
	}
	if cfg.Database.DBName == "" {
		errors = append(errors, "Database name is not configured (set DATABASE_URL or database.dbname in config.yaml)")
	}

	// 验证 JWT 密钥
	if cfg.Auth.JWTSecret == "" || cfg.Auth.JWTSecret == "change-me" {
		errors = append(errors, "JWT secret is not configured or using default value (set JWT_SECRET environment variable or auth.jwt_secret in config.yaml)")
	}

	// 验证 gRPC 地址
	if cfg.GRPC.AIServiceAddr == "" {
		errors = append(errors, "gRPC AI service address is not configured (set GRPC_AI_SERVICE_ADDR or grpc.ai_service_addr in config.yaml)")
	}

	if len(errors) > 0 {
		log.Println("Configuration validation failed:")
		for _, err := range errors {
			log.Printf("  - %s", err)
		}
		os.Exit(1)
	}
}

// logConfigSources 输出配置来源日志
func logConfigSources(cfg *Config) {
	log.Println("Configuration loaded successfully:")
	for key, source := range cfg.configSources {
		log.Printf("  - %s: %s", key, source)
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
