package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"strings"
)

func (c *Client) ListContainers() ([][]string, error) {
	containers, err := c.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result [][]string
	for _, container := range containers {
		result = append(result, []string{
			container.ID[:12],
			container.Image,
			truncate(container.Command, 20),
			formatTime(container.Created),
			container.Status,
			formatPorts(container.Ports),
			strings.TrimPrefix(container.Names[0], "/"),
		})
	}
	return result, nil
}

func (c *Client) GetContainerLogs(containerID string) (string, error) {
	options := container.LogsOptions{ShowStdout: true, ShowStderr: true}
	logs, err := c.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		return "", err
	}
	defer logs.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(logs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Client) InspectContainer(containerID string) (string, error) {
	_, rawJSON, err := c.ContainerInspectWithRaw(context.Background(), containerID, false)
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

func (c *Client) StartContainer(containerID string) error {
	return c.ContainerStart(context.Background(), containerID, container.StartOptions{})
}

func (c *Client) StopContainer(containerID string) error {
	timeout := 10
	stopOptions := container.StopOptions{Timeout: &timeout}
	return c.ContainerStop(context.Background(), containerID, stopOptions)
}

func (c *Client) PauseContainer(containerID string) error {
	return c.ContainerPause(context.Background(), containerID)
}

func (c *Client) UnpauseContainer(containerID string) error {
	return c.ContainerUnpause(context.Background(), containerID)
}

func (c *Client) KillContainer(containerID string) error {
	return c.ContainerKill(context.Background(), containerID, "SIGKILL")
}

func (c *Client) DeleteContainer(containerID string) error {
	return c.ContainerRemove(context.Background(), containerID, container.RemoveOptions{})
}
