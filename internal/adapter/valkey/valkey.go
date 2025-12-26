package valkey

import (
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Client struct {
	valkey.Client

	ttl time.Duration
}

func New(cfg Config) (*Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: cfg.Address,
		Username:    cfg.Username,
		Password:    cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("valkey: %w", err)
	}
	return &Client{
		Client: client,
		ttl:    cfg.TTL,
	}, nil
}

func MustNew(cfg Config) *Client {
	client, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func (c *Client) settingKey(userID, chatID int64) string {
	return fmt.Sprintf("settings:%d:%d", userID, chatID)
}
