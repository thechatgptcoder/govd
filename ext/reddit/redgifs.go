package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/thechatgptcoder/govd/enums"
	"github.com/thechatgptcoder/govd/models"
)

var redgifsAPI = regexp.MustCompile(`https://www\.redgifs\.com/watch/([a-zA-Z0-9]+)`)

type redgifsResponse struct {
	Gif struct {
		VideoURL string `json:"urls"`
	} `json:"gif"`
}

func ExtractRedgifs(ctx context.Context, url string, task *models.Task) error {
	// Extract Redgifs ID
	matches := redgifsAPI.FindStringSubmatch(url)
	if len(matches) < 2 {
		return fmt.Errorf("invalid redgifs URL: %s", url)
	}
	id := matches[1]

	// Query Redgifs API (no auth needed for public videos)
	apiURL := fmt.Sprintf("https://api.redgifs.com/v2/gifs/%s", id)
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	req.Header.Set("User-Agent", "govd-bot")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching redgifs API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("redgifs API returned status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Gif struct {
			MP4URL string `json:"urls"`
		} `json:"gif"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("failed to parse redgifs JSON: %w", err)
	}

	if data.Gif.MP4URL == "" {
		return fmt.Errorf("no video URL found in Redgifs response")
	}

	media := &models.Media{
		Caption: task.Title,
		Formats: []*models.MediaFormat{
			{
				FormatID:   "mp4",
				Type:       enums.MediaTypeVideo,
				VideoCodec: enums.MediaCodecAVC,
				URL:        []string{data.Gif.MP4URL},
			},
		},
	}
	task.MediaList = []*models.Media{media}
	return nil
}
