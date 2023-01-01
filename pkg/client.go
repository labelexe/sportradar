package sportradar

import (
	"github.com/playback-sports/sportradar/pkg/nba"
	"github.com/playback-sports/sportradar/pkg/ncaaf"
	"github.com/playback-sports/sportradar/pkg/ncaamb"
	"github.com/playback-sports/sportradar/pkg/nhl"
	"net/http"
	"time"

	sr "github.com/playback-sports/sportradar/pkg/base"
	"github.com/playback-sports/sportradar/pkg/mlb"
	"github.com/playback-sports/sportradar/pkg/nfl"
)

type LeagueKeys struct {
	MLB    string `json:"mlb"`
	NFL    string `json:"nfl"`
	NBA    string `json:"nba"`
	NHL    string `json:"nhl"`
	NCAAF  string `json:"ncaaf"`
	NCAAMB string `json:"ncaamb"`
}

type Client struct {
	c   http.Client
	cfg ClientConfig
}

type ClientConfig struct {
	Keys LeagueKeys
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

func (c Client) MLBSchedule(year int, st sr.SeasonType) (mlb.Schedule, error) {
	return mlb.FetchSchedule(c.c, year, st, c.cfg.Keys.MLB)
}

func (c Client) NFLSchedule(year int, st sr.SeasonType) (nfl.Schedule, error) {
	return nfl.FetchSchedule(c.c, year, st, c.cfg.Keys.NFL)
}

func (c Client) NBASchedule(year int, st sr.SeasonType) (nba.Schedule, error) {
	return nba.FetchSchedule(c.c, year, st, c.cfg.Keys.NBA)
}

func (c Client) NHLSchedule(year int, st sr.SeasonType) (nhl.Schedule, error) {
	return nhl.FetchSchedule(c.c, year, st, c.cfg.Keys.NHL)
}

func (c Client) NCAAFSchedule(year int, st sr.SeasonType) (ncaaf.Schedule, error) {
	return ncaaf.FetchSchedule(c.c, year, st, c.cfg.Keys.NCAAF)
}

func (c Client) NCAAMBSchedule(year int, st sr.SeasonType) (ncaamb.Schedule, error) {
	return ncaamb.FetchSchedule(c.c, year, st, c.cfg.Keys.NCAAMB)
}
