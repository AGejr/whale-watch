package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"strings"
	"time"
)

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatPorts(ports []types.Port) string {
	var result []string
	for _, port := range ports {
		result = append(result, fmt.Sprintf("%d/%s", port.PrivatePort, port.Type))
	}
	return strings.Join(result, ", ")
}

func formatTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

func parseRepoTag(repoTag string) (string, string) {
	parts := strings.Split(repoTag, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return repoTag, "<none>"
}
