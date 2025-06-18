package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool     *pgxpool.Pool
	poolConn *pgxpool.Conn
	conn     *pgx.Conn
}

type Options struct {
	ConnectionURI   string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	QueryTimeout    time.Duration
}

const (
	DefaultQueryTimeout    = 10 * time.Second
	DefaultMaxOpenConns    = 25
	DefaultMaxIdleConns    = 25
	DefaultConnMaxLifetime = 5 * time.Minute
)

func DefaultOptions() *Options {
	return &Options{
		MaxOpenConns:    DefaultMaxOpenConns,
		MaxIdleConns:    DefaultMaxIdleConns,
		ConnMaxLifetime: DefaultConnMaxLifetime,
		QueryTimeout:    DefaultQueryTimeout,
	}
}

func NewService(ctx context.Context, opts *Options) (*Service, error) {
	if opts == nil || opts.ConnectionURI == "" {
		return nil, fmt.Errorf("connection uri is required")
	}

	pool, err := pgxpool.New(ctx, opts.ConnectionURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	pool.Config().MaxConns = int32(opts.MaxOpenConns)
	pool.Config().MaxConnIdleTime = opts.ConnMaxLifetime

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Service{
		pool: pool,
	}, nil
}

func (s *Service) PoolConn(ctx context.Context) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire database connection: %w", err)
	}

	s.poolConn = conn
	return nil
}

func (s *Service) Conn() *pgx.Conn {
	if s.conn != nil {
		return s.conn
	}

	return s.poolConn.Conn()
}

func (s *Service) Disconnect(ctx context.Context) {
	if s.conn != nil {
		s.conn.Close(ctx)
	}

	if s.poolConn != nil {
		s.poolConn.Release()
	}

	if s.pool != nil {
		s.pool.Close()
	}
}
