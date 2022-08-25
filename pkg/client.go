package sportradar

import (
	"net/http"
	"time"

	sr "github.com/playback-sports/sportradar/pkg/base"
	"github.com/playback-sports/sportradar/pkg/mlb"
	"github.com/playback-sports/sportradar/pkg/nfl"
)

type LeageuKeys struct {
	MLB string
	NFL string
}

type Client struct {
	c   http.Client
	cfg ClientConfig
}

type ClientConfig struct {
	Keys LeageuKeys
}

func NewClient(cfg ClientConfig) Client {
	return Client{
		c: http.Client{
			Timeout: 5 * time.Minute,
		},
		cfg: cfg,
	}
}

func (c Client) MLBDailySummary(t time.Time) (mlb.DailySummary, error) {
	return mlb.FetchDailySummary(c.c, t, c.cfg.Keys.MLB)
}

func (c Client) MLBSchedule(t time.Time, st sr.SeasonType) (mlb.Schedule, error) {
	return mlb.FetchSchedule(c.c, t, st, c.cfg.Keys.MLB)
}

func (c Client) NFLSchedule(t time.Time, st sr.SeasonType) (nfl.Schedule, error) {
	return nfl.FetchSchedule(c.c, t, st, c.cfg.Keys.NFL)
}
