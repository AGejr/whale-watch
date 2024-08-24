package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types/image"
)

func (c *Client) ListImages() ([][]string, error) {
	images, err := c.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result [][]string
	for _, image := range images {
		repoTags := "<none>:<none>"
		if len(image.RepoTags) > 0 {
			repoTags = image.RepoTags[0]
		}
		repo, tag := parseRepoTag(repoTags)
		result = append(result, []string{
			repo,
			tag,
			image.ID[7:19],
			formatTime(image.Created),
			formatSize(image.Size),
		})
	}
	return result, nil
}

func (c *Client) InspectImage(imageID string) (string, error) {
	_, rawJSON, err := c.ImageInspectWithRaw(context.Background(), imageID)
	if err != nil {
		return "", err
	}
	var formattedJSON bytes.Buffer
	err = json.Indent(&formattedJSON, rawJSON, "", "\t")
	if err != nil {
		return "", err
	}
	return string(formattedJSON.Bytes()), nil
}

func (c *Client) DeleteImage(imageID string) error {
	_, err := c.ImageRemove(context.Background(), imageID, image.RemoveOptions{Force: false, PruneChildren: true})
	return err
}
