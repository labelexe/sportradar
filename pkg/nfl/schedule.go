package nfl

import (
	"encoding/json"
	"fmt"
	sr "github.com/playback-sports/sportradar/pkg/base"
	"io/ioutil"
	"net/http"
)

type srScheduleBroadcast struct {
	Network string `json:"network"`
}

type srScheduleTeam struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Alias      string `json:"alias"`
	GameNumber int    `json:"game_number"`
	SrID       string `json:"sr_id"`
}

type srScheduleGame struct {
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

type srScheduleWeek struct {
	ID       string           `json:"id"`
	Sequence int              `json:"sequence"`
	Title    string           `json:"title"`
	Games    []srScheduleGame `json:"games"`
}

type srSchedule struct {
	ID      string           `json:"id"`
	Year    int              `json:"year"`
	Type    string           `json:"type"`
	Name    string           `json:"name"`
	Weeks   []srScheduleWeek `json:"weeks"`
	Comment string           `json:"_comment"`
}

type ScheduleGame struct {
	ID        string              `json:"id"`
	Status    sr.GameStatus       `json:"status"`
	Scheduled string              `json:"scheduled"`
	EntryMode string              `json:"entry_mode"`
	SrID      string              `json:"sr_id"`
	GameType  string              `json:"game_type"`
	Venue     sr.Venue            `json:"venue"`
	HomeTeam  srScheduleTeam      `json:"home"`
	AwayTeam  srScheduleTeam      `json:"away"`
	Broadcast srScheduleBroadcast `json:"broadcast"`
}

func srScheduleGameConvert(g srScheduleGame) ScheduleGame {
	return ScheduleGame{
		ID:        g.ID,
		Status:    sr.ParseGameStatus(g.Status),
		Scheduled: g.Scheduled,
		EntryMode: g.EntryMode,
		SrID:      g.SrID,
		GameType:  g.GameType,
		Venue:     g.Venue,
		HomeTeam:  g.HomeTeam,
		AwayTeam:  g.AwayTeam,
		Broadcast: g.Broadcast,
	}
}

type ScheduleWeek struct {
	ID       string         `json:"id"`
	Sequence int            `json:"sequence"`
	Title    string         `json:"title"`
	Games    []ScheduleGame `json:"games"`
}

func srScheduleWeekConvert(w srScheduleWeek) ScheduleWeek {
	games := make([]ScheduleGame, 0, len(w.Games))
	for _, game := range w.Games {
		games = append(games, srScheduleGameConvert(game))
	}
	return ScheduleWeek{
		ID:       w.ID,
		Sequence: w.Sequence,
		Title:    w.Title,
		Games:    games,
	}
}

type Schedule struct {
	ID      string         `json:"id"`
	Year    int            `json:"year"`
	Type    sr.SeasonType  `json:"type"`
	Name    string         `json:"name"`
	Weeks   []ScheduleWeek `json:"weeks"`
	Comment string         `json:"_comment"`
}

func srScheduleConvert(s srSchedule) Schedule {
	weeks := make([]ScheduleWeek, 0, len(s.Weeks))
	for _, week := range s.Weeks {
		weeks = append(weeks, srScheduleWeekConvert(week))
	}
	return Schedule{
		ID:      s.ID,
		Year:    s.Year,
		Type:    sr.ParseSeasonType(s.Type),
		Name:    s.Name,
		Weeks:   weeks,
		Comment: s.Comment,
	}
}

var scheduleURLTemplate = "https://api.sportradar.us/nfl/official/production/v7/en/games/%d/%s/schedule.json?api_key=%s"

func scheduleURL(year int, st sr.SeasonType, apiKey string) string {
	return fmt.Sprintf(scheduleURLTemplate, year, st, apiKey)
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func FetchSchedule(c http.Client, year int, st sr.SeasonType, apiKey string) (Schedule, error) {
	url := scheduleURL(year, st, apiKey)
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
		return Schedule{}, fmt.Errorf("invalid response code from sportradar nfl schedule request: %s %d", url, resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Schedule{}, err
	}

	schedule := srSchedule{}
	err = json.Unmarshal(respBytes, &schedule)
	if err != nil {
		return Schedule{}, err
	}

	return srScheduleConvert(schedule), nil
}
