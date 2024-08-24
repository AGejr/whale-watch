package docker

import (
	"context"
	"github.com/docker/docker/api/types/volume"
)

func (c *Client) ListVolumes() ([][]string, error) {
	volumes, err := c.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result [][]string
	for _, volume := range volumes.Volumes {
		result = append(result, []string{
			volume.Driver,
			volume.Name,
			volume.Mountpoint,
		})
	}
	return result, nil
}

func (c *Client) DeleteVolume(volumeName string) error {
	return c.VolumeRemove(context.Background(), volumeName, false)
}
