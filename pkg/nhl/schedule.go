package nhl

import (
	"encoding/json"
	"fmt"
	sr "github.com/playback-sports/sportradar/pkg/base"
	"io/ioutil"
	"net/http"
	"time"
)

var scheduleURLTemplate = "https://api.sportradar.us/nhl/production/v7/en/games/%d/%s/schedule.json?api_key=%s"

func scheduleURL(t time.Time, st sr.SeasonType, apiKey string) string {
	return fmt.Sprintf(scheduleURLTemplate, t.Year(), st, apiKey)
}

type Schedule struct {
	League srLeague       `json:"league"`
	Season srSeason       `json:"season"`
	Games  []ScheduleGame `json:"games"`
}

type srLeague struct {
	ID string `json:"id"`
}

type srSeason struct {
	ID   string `json:"id"`
	Year int    `json:"year"`
	Type string `json:"type"`
}

type ScheduleGame struct {
	ID         string                `json:"id"`
	Status     sr.GameStatus         `json:"status"`
	Scheduled  string                `json:"scheduled"`
	SrID       string                `json:"sr_id"`
	GameType   string                `json:"game_type"`
	Venue      sr.Venue              `json:"venue"`
	HomeTeam   srScheduleTeam        `json:"home"`
	AwayTeam   srScheduleTeam        `json:"away"`
	Broadcasts []srScheduleBroadcast `json:"broadcasts"`
}

type srScheduleBroadcast struct {
	Network string `json:"network"`
	Type    string `json:"type"`
	Locale  string `json:"locale"`
	Channel string `json:"channel"`
}

type srScheduleTeam struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Alias     string `json:"alias"`
	Reference string `json:"reference"`
	SrID      string `json:"sr_id"`
}

func FetchSchedule(c http.Client, t time.Time, st sr.SeasonType, apiKey string) (Schedule, error) {
	url := scheduleURL(t, st, apiKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Schedule{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return Schedule{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Schedule{}, fmt.Errorf("invalid response code from sportradar nba schedule request: %s %d", url, resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Schedule{}, err
	}

	schedule := Schedule{}
	err = json.Unmarshal(respBytes, &schedule)
	if err != nil {
		return Schedule{}, err
	}

	return schedule, nil
}
