package plugins

import (
	"strings"

	"github.com/govdbot/govd/extractor"
	"github.com/govdbot/govd/models"
	"github.com/govdbot/govd/plugin"
)

func init() {
	plugin.Register("reddit-redgifs", plugin.Rule{
		Matches: func(url string) bool {
			// Detect Reddit posts that embed redgifs
			return strings.Contains(url, "reddit.com") && strings.Contains(url, "redgifs.com")
		},
		Process: func(task *models.Task) error {
			// Extract the redgifs URL from the Reddit post
			if task.Reddit != nil && task.Reddit.EmbeddedURL != "" && strings.Contains(task.Reddit.EmbeddedURL, "redgifs.com") {
				// Delegate to the redgifs extractor
				return extractor.Extract(task.Reddit.EmbeddedURL, task)
			}
			return nil
		},
	})
}
