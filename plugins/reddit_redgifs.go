package plugins

import (
	"context"
	"strings"

	"github.com/thechatgptcoder/govd/ext/reddit"
	"github.com/thechatgptcoder/govd/models"
	"github.com/thechatgptcoder/govd/plugin"
)

func init() {
	plugin.Register("reddit-redgifs", plugin.Rule{
		Matches: func(url string) bool {
			// Detect Reddit posts that embed redgifs
			return strings.Contains(url, "reddit.com") && strings.Contains(url, "redgifs.com")
		},
		Process: func(ctx context.Context, task *models.Task) error {
			if task.Reddit != nil && task.Reddit.EmbeddedURL != "" && strings.Contains(task.Reddit.EmbeddedURL, "redgifs.com") {
				return reddit.ExtractRedgifs(ctx, task.Reddit.EmbeddedURL, task)
			}
			return nil
		},
	})
}

