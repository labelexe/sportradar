package mlb

import (
	"encoding/json"
	"fmt"
	sr "github.com/playback-sports/sportradar/pkg/base"
	"io/ioutil"
	"net/http"
)

type srScheduleLeague struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
	ID    string `json:"id"`
}

type srScheduleSeason struct {
	ID   string `json:"id"`
	Year int    `json:"year"`
	Type string `json:"type"`
}

type srScheduleTeam struct {
	Name   string `json:"name"`
	Market string `json:"market"`
	Abbr   string `json:"abbr"`
	ID     string `json:"id"`
}

type srScheduleBroadcast struct {
	Network string `json:"network"`
}

type srScheduleGame struct {
	ID           string              `json:"id"`
	Status       string              `json:"status"`
	Coverage     string              `json:"coverage"`
	GameNumber   int                 `json:"game_number"`
	DayNight     string              `json:"day_night"`
	Scheduled    string              `json:"scheduled"`
	HomeTeamID   string              `json:"home_team"`
	AwayTeamID   string              `json:"away_team"`
	Attendance   int                 `json:"attendance"`
	Duration     string              `json:"duration"`
	DoubleHeader bool                `json:"double_header"`
	EntryMode    string              `json:"entry_mode"`
	Reference    string              `json:"reference"`
	Venue        sr.Venue            `json:"venue"`
	HomeTeam     srScheduleTeam      `json:"home"`
	AwayTeam     srScheduleTeam      `json:"away"`
	Broadcast    srScheduleBroadcast `json:"broadcast"`
}

type srSchedule struct {
	League srScheduleLeague `json:"league"`
	Season srScheduleSeason `json:"season"`
	Games  []srScheduleGame `json:"games"`
}

type ScheduleGame struct {
	ID           string              `json:"id"`
	Status       sr.GameStatus       `json:"status"`
	Coverage     string              `json:"coverage"`
	GameNumber   int                 `json:"game_number"`
	DayNight     string              `json:"day_night"`
	Scheduled    string              `json:"scheduled"`
	HomeTeamID   string              `json:"home_team"`
	AwayTeamID   string              `json:"away_team"`
	Attendance   int                 `json:"attendance"`
	Duration     string              `json:"duration"`
	DoubleHeader bool                `json:"double_header"`
	EntryMode    string              `json:"entry_mode"`
	Reference    string              `json:"reference"`
	Venue        sr.Venue            `json:"venue"`
	HomeTeam     srScheduleTeam      `json:"home"`
	AwayTeam     srScheduleTeam      `json:"away"`
	Broadcast    srScheduleBroadcast `json:"broadcast"`
}

func srScheduleGameConvert(game srScheduleGame) ScheduleGame {
	return ScheduleGame{
		ID:           game.ID,
		Status:       sr.ParseGameStatus(game.Status),
		Coverage:     game.Coverage,
		GameNumber:   game.GameNumber,
		DayNight:     game.DayNight,
		Scheduled:    game.Scheduled,
		HomeTeamID:   game.HomeTeamID,
		AwayTeamID:   game.AwayTeamID,
		Attendance:   game.Attendance,
		Duration:     game.Duration,
		DoubleHeader: game.DoubleHeader,
		EntryMode:    game.EntryMode,
		Reference:    game.Reference,
		Venue:        game.Venue,
		HomeTeam:     game.HomeTeam,
		AwayTeam:     game.AwayTeam,
		Broadcast:    game.Broadcast,
	}
}

type Schedule struct {
	League srScheduleLeague `json:"league"`
	Season srScheduleSeason `json:"season"`
	Games  []ScheduleGame   `json:"games"`
}

func srScheduleConvert(schedule srSchedule) Schedule {
	games := make([]ScheduleGame, 0, len(schedule.Games))
	for _, game := range schedule.Games {
		games = append(games, srScheduleGameConvert(game))
	}
	return Schedule{
		League: schedule.League,
		Season: schedule.Season,
		Games:  games,
	}
}

var scheduleURLTemplate = "https://api.sportradar.com/mlb/production/v7/en/games/%d/%s/schedule.json?api_key=%s"

func scheduleURL(year int, st sr.SeasonType, apiKey string) string {
	return fmt.Sprintf(scheduleURLTemplate, year, st, apiKey)
}

func FetchSchedule(c http.Client, year int, st sr.SeasonType, apiKey string) (Schedule, error) {
	url := scheduleURL(year, st, apiKey)
	// fmt.Printf("schedule url: %s", url)
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
		return Schedule{}, fmt.Errorf("invalid response code from sportradar mlb schedule request: %s %d", url, resp.StatusCode)
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
