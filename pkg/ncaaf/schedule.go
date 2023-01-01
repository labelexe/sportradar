package ncaaf

import (
	"encoding/json"
	"fmt"
	sr "github.com/playback-sports/sportradar/pkg/base"
	"io/ioutil"
	"net/http"
	"time"
)

var scheduleURLTemplate = "https://api.sportradar.us/ncaafb/production/v7/en/games/%d/%s/schedule.json?api_key=%s"

func scheduleURL(t time.Time, st sr.SeasonType, apiKey string) string {
	return fmt.Sprintf(scheduleURLTemplate, t.Year(), st, apiKey)
}

type srScheduleBroadcast struct {
	Network string `json:"network"`
}

type srScheduleTeam struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Alias      string `json:"alias"`
	GameNumber int    `json:"game_number"`
}

type ScheduleGame struct {
	ID        string              `json:"id"`
	Status    string              `json:"status"`
	Scheduled string              `json:"scheduled"`
	EntryMode string              `json:"entry_mode"`
	SrID      string              `json:"sr_id"`
	GameType  string              `json:"game_type"`
	Venue     sr.Venue            `json:"venue"`
	HomeTeam  srScheduleTeam      `json:"home"`
	AwayTeam  srScheduleTeam      `json:"away"`
	Broadcast srScheduleBroadcast `json:"broadcast"`
}

type ScheduleWeek struct {
	ID       string         `json:"id"`
	Sequence int            `json:"sequence"`
	Title    string         `json:"title"`
	Games    []ScheduleGame `json:"games"`
}

type Schedule struct {
	ID      string         `json:"id"`
	Year    int            `json:"year"`
	Type    sr.SeasonType  `json:"type"`
	Name    string         `json:"name"`
	Weeks   []ScheduleWeek `json:"weeks"`
	Comment string         `json:"_comment"`
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
		return Schedule{}, fmt.Errorf("invalid response code from sportradar ncaaf schedule request: %s %d", url, resp.StatusCode)
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
