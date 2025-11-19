package valkey

import "github.com/valkey-io/valkey-go"

const (
	summaryKeyPrefix = "summary"
)

type Config struct {
	URL string `env:"VALKEY_URL,required"`
}

type Client struct {
	client valkey.Client
}

func New(cfg Config) *Client {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{cfg.URL},
	})
	if err != nil {
		panic(err)
	}
	return &Client{
		client: client,
	}
}
