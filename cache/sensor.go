package cache

type Sensor struct {
	Uuid   string `json:"uuid"`
	Active bool   `json:"active"`
	Type   string `json:"type"`
	Max    string `json:"max"`
	Min    string `json:"min"`
}
