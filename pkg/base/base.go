package base

type SeasonType string

const (
	SeasonTypePre     SeasonType = "PRE"
	SeasonTypeRegular SeasonType = "Reg"
	SeasonTypePost    SeasonType = "PST"
)

type GameStatus string

const (
	GameStatusScheduled  GameStatus = "scheduled"
	GameStatusClosed     GameStatus = "closed"
	GameStatusInProgress GameStatus = "inprogress"
)

func ParseGameStatus(s string) GameStatus {
	if s == string(GameStatusScheduled) {
		return GameStatusScheduled
	} else if s == string(GameStatusInProgress) {
		return GameStatusInProgress
	} else {
		return GameStatusClosed
	}
}

type Location struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lng"`
}

type Weather struct {
	Forecast Forecast `json:"forecast"`
}

type Forecast struct {
	TempF      int    `json:"temp_f"`
	Condition  string `json:"condition"`
	Humidity   int    `json:"humidity"`
	DewPointF  int    `json:"dew_point_f"`
	CloudCover int    `json:"cloud_cover"`
	ObsTime    string `json:"obs_time"`
	Wind       Wind   `jsong:"wind"`
}

type Wind struct {
	SpeedMPH  int    `json:"speed_mph"`
	Direction string `json:"direction"`
}

type Venue struct {
	Name             string   `json:"name"`
	Market           string   `json:"market"`
	Capacity         int      `json:"capacity"`
	Surface          string   `json:"surface"`
	Address          string   `json:"address"`
	City             string   `json:"city"`
	State            string   `json:"state"`
	Zip              string   `json:"zip"`
	Country          string   `json:"country"`
	ID               string   `json:"id"`
	FieldOrientation string   `json:"field_orientation"`
	StadiumType      string   `json:"stadium_type"`
	Location         Location `json:"location"`
}

type Broadcast struct {
	Network string `json:"network"`
}
