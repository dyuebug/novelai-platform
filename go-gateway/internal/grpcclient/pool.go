package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"

	"go-gateway/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnPool gRPC 连接池
type ConnPool struct {
	addr     string
	maxConns int
	conns    chan *grpc.ClientConn
	mu       sync.Mutex
	tlsCfg   *tls.Config
}

// NewConnPool 创建连接池
func NewConnPool(addr string, maxConns int) *ConnPool {
	if maxConns <= 0 {
		maxConns = 4
	}
	return &ConnPool{
		addr:     addr,
		maxConns: maxConns,
		conns:    make(chan *grpc.ClientConn, maxConns),
	}
}

// NewConnPoolWithTLS 创建支持 mTLS 的连接池
func NewConnPoolWithTLS(cfg config.GRPCConfig) (*ConnPool, error) {
	pool := &ConnPool{
		addr:     cfg.AIServiceAddr,
		maxConns: cfg.MaxConns,
		conns:    make(chan *grpc.ClientConn, cfg.MaxConns),
	}

	if cfg.TLSEnabled {
		tlsCfg, err := loadTLSConfig(cfg.CertFile, cfg.KeyFile, cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}
		pool.tlsCfg = tlsCfg
	}

	return pool, nil
}

// loadTLSConfig 加载 mTLS 配置
func loadTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	// 加载客户端证书
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// 加载 CA 证书
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// Acquire 获取连接
func (p *ConnPool) Acquire(ctx context.Context) (*grpc.ClientConn, func(), error) {
	select {
	case conn := <-p.conns:
		return conn, func() { p.release(conn) }, nil
	default:
		conn, err := p.createConn(ctx)
		if err != nil {
			return nil, func() {}, err
		}
		return conn, func() { p.release(conn) }, nil
	}
}

// createConn 创建新连接
func (p *ConnPool) createConn(ctx context.Context) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if p.tlsCfg != nil {
		// 使用 mTLS
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(p.tlsCfg)))
	} else {
		// 不安全连接 (开发环境)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	opts = append(opts, grpc.WithDefaultCallOptions(grpc.CallContentSubtype(jsonCodecName)))

	return grpc.DialContext(ctx, p.addr, opts...)
}

func (p *ConnPool) release(conn *grpc.ClientConn) {
	select {
	case p.conns <- conn:
		// 放回池中
	default:
		// 池满，关闭连接
		_ = conn.Close()
	}
}

// Close 关闭连接池
func (p *ConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	close(p.conns)
	for conn := range p.conns {
		_ = conn.Close()
	}
}

// IsTLSEnabled 检查是否启用 TLS
func (p *ConnPool) IsTLSEnabled() bool {
	return p.tlsCfg != nil
}
