package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

type Config struct {
	URL                   string `env:"VALKEY_URL,required"`
	RingScaleEachConn     int    `env:"VALKEY_RING_SCALE_EACH_CONN" envDefault:"10"`
	BlockingPoolSize      int    `env:"VALKEY_BLOCKING_POOL_SIZE" envDefault:"1024"`
	BlockingPipeline      int    `env:"VALKEY_BLOCKING_PIPELINE" envDefault:"2000"`
	ReadBufferEachConn    int    `env:"VALKEY_READ_BUFFER_EACH_CONN" envDefault:"524288"`   // 0.5 MiB
	WriteBufferEachConn   int    `env:"VALKEY_WRITE_BUFFER_EACH_CONN" envDefault:"524288"`  // 0.5 MiB
	CacheSizeEachConn     int    `env:"VALKEY_CACHE_SIZE_EACH_CONN" envDefault:"134217728"` // 128 MiB
	DisableCache          bool   `env:"VALKEY_DISABLE_CACHE" envDefault:"false"`
	DisableRetry          bool   `env:"VALKEY_DISABLE_RETRY" envDefault:"false"`
	AlwaysPipelining      bool   `env:"VALKEY_ALWAYS_PIPELINING" envDefault:"false"`
	DisableAutoPipelining bool   `env:"VALKEY_DISABLE_AUTO_PIPELINING" envDefault:"false"`
}

type Client struct {
	client valkey.Client
}

func New(cfg Config) *Client {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress:           []string{cfg.URL},
		RingScaleEachConn:     cfg.RingScaleEachConn,
		BlockingPoolSize:      cfg.BlockingPoolSize,
		BlockingPipeline:      cfg.BlockingPipeline,
		ReadBufferEachConn:    cfg.ReadBufferEachConn,
		WriteBufferEachConn:   cfg.WriteBufferEachConn,
		CacheSizeEachConn:     cfg.CacheSizeEachConn,
		DisableCache:          cfg.DisableCache,
		DisableRetry:          cfg.DisableRetry,
		AlwaysPipelining:      cfg.AlwaysPipelining,
		DisableAutoPipelining: cfg.DisableAutoPipelining,
	})
	if err != nil {
		panic(err)
	}
	cmd := client.Do(context.Background(), client.B().Ping().Build())
	if err := cmd.Error(); err != nil {
		panic("valkey ping failed: " + err.Error())
	}
	return &Client{
		client: client,
	}
}

func (c *Client) Close() error {
	c.client.Close()
	return nil
}

func (c *Client) summaryKey(ownerID, peerID string) string {
	return fmt.Sprintf("summary:%s:%s", ownerID, peerID)
}
