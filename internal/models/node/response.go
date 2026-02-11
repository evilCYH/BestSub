package node

type Response struct {
	SubID       uint16 `json:"sub_id"`
	UniqueKey   uint64 `json:"unique_key"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Delay       uint16 `json:"delay"`
	SpeedUp     uint32 `json:"speed_up"`
	SpeedDown   uint32 `json:"speed_down"`
	Risk        uint8  `json:"risk"`
	AliveStatus uint64 `json:"alive_status"`
	Country     string `json:"country"`
}

