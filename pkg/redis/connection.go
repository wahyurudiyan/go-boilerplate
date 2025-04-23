package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/wahyurudiyan/go-boilerplate/pkg/config"

	goRedis "github.com/redis/go-redis/v9"
)

type redisClient struct {
	cfg *RedisConfig
}

func NewClient() (*goRedis.Client, error) {
	var cfg RedisConfig
	config.Load(&cfg)

	newCli := redisClient{
		cfg: &cfg,
	}

	client, err := newCli.connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *redisClient) connect() (*goRedis.Client, error) {
	if !r.cfg.RedisIsCluster {
		return r.connectRedisStandaloneClient()
	}
	return r.connectRedisStandaloneClient()
}

// BuildRedisStandaloneClient creates a standalone Redis client
func (r *redisClient) connectRedisStandaloneClient() (*goRedis.Client, error) {
	options := &goRedis.Options{
		Addr:     r.cfg.RedisAddr,
		Password: r.cfg.RedisPassword,
		DB:       r.cfg.RedisDB,

		PoolSize:        r.cfg.RedisPoolSize,
		MinIdleConns:    r.cfg.RedisMinIdleConns,
		PoolTimeout:     r.cfg.RedisPoolTimeout,
		MaxRetries:      r.cfg.RedisMaxRetries,
		MinRetryBackoff: r.cfg.RedisMinRetryBackoff,
		MaxRetryBackoff: r.cfg.RedisMaxRetryBackoff,

		DialTimeout:     r.cfg.RedisDialTimeout,
		ReadTimeout:     r.cfg.RedisReadTimeout,
		WriteTimeout:    r.cfg.RedisWriteTimeout,
		ConnMaxIdleTime: r.cfg.RedisIdleTimeout,
		OnConnect: func(ctx context.Context, cn *goRedis.Conn) error {
			_, err := cn.Ping(ctx).Result()
			return err
		},
	}

	if r.cfg.RedisEnableTLS {
		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: r.cfg.RedisTLSInsecureSkipVerify,
		}
	}

	client := goRedis.NewClient(options)
	// Perform a test connection
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.RedisDialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

func (c *RedisConfig) connectRedisClusterClient() (*goRedis.ClusterClient, error) {
	if len(c.RedisClusterAddrs) == 0 {
		return nil, fmt.Errorf("no cluster addresses provided for Redis cluster")
	}

	options := &goRedis.ClusterOptions{
		Addrs:    c.RedisClusterAddrs,
		Password: c.RedisPassword,

		PoolSize:        c.RedisPoolSize,
		MinIdleConns:    c.RedisMinIdleConns,
		PoolTimeout:     c.RedisPoolTimeout,
		MaxRetries:      c.RedisMaxRetries,
		MinRetryBackoff: c.RedisMinRetryBackoff,
		MaxRetryBackoff: c.RedisMaxRetryBackoff,

		DialTimeout:     c.RedisDialTimeout,
		ReadTimeout:     c.RedisReadTimeout,
		WriteTimeout:    c.RedisWriteTimeout,
		ConnMaxIdleTime: c.RedisIdleTimeout,

		MaxRedirects:   c.RedisMaxRedirects,
		RouteByLatency: c.RedisRouteByLatency,
		RouteRandomly:  c.RedisRouteRandomly,
		ReadOnly:       c.RedisReadOnly,

		OnConnect: func(ctx context.Context, cn *goRedis.Conn) error {
			_, err := cn.Ping(ctx).Result()
			return err
		},
	}

	if c.RedisEnableTLS {
		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: c.RedisTLSInsecureSkipVerify,
		}
	}

	client := goRedis.NewClusterClient(options)
	// Perform a test connection
	ctx, cancel := context.WithTimeout(context.Background(), c.RedisDialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis cluster: %w", err)
	}

	return client, nil
}
