// Package client provides a wrapper around the Zenon SDK RPC client
// with connection management and CLI-specific helpers.
package client

import (
	"fmt"
	"time"

	"github.com/0x3639/znn-sdk-go/rpc_client"
)

// Client wraps the SDK RpcClient with CLI-specific functionality
type Client struct {
	*rpc_client.RpcClient
	url string
}

// New creates a new RPC client with the specified URL and default options.
// The client will automatically reconnect on connection loss.
func New(url string) (*Client, error) {
	if url == "" {
		url = "ws://127.0.0.1:35998"
	}

	opts := rpc_client.DefaultClientOptions()
	opts.AutoReconnect = true
	opts.ReconnectDelay = 2 * time.Second
	opts.MaxReconnectDelay = 60 * time.Second
	opts.ReconnectAttempts = 10
	opts.HealthCheckInterval = 15 * time.Second

	client, err := rpc_client.NewRpcClientWithOptions(url, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node at %s: %w", url, err)
	}

	return &Client{
		RpcClient: client,
		url:       url,
	}, nil
}

// NewWithOptions creates a new RPC client with custom options
func NewWithOptions(url string, opts rpc_client.ClientOptions) (*Client, error) {
	if url == "" {
		url = "ws://127.0.0.1:35998"
	}

	client, err := rpc_client.NewRpcClientWithOptions(url, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node at %s: %w", url, err)
	}

	return &Client{
		RpcClient: client,
		url:       url,
	}, nil
}

// URL returns the WebSocket URL this client is connected to
func (c *Client) URL() string {
	return c.url
}

// Close stops the client and closes the connection
func (c *Client) Close() error {
	c.Stop()
	return nil
}
