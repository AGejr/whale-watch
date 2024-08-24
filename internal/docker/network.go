package docker

import (
	"context"
	"github.com/docker/docker/api/types/network"
)

func (c *Client) ListNetworks() ([][]string, error) {
	networks, err := c.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result [][]string
	for _, network := range networks {
		result = append(result, []string{
			network.ID,
			network.Name,
			network.Driver,
			network.Scope,
		})
	}
	return result, nil
}

func (c *Client) DeleteNetwork(networkID string) error {
	return c.NetworkRemove(context.Background(), networkID)
}
