package node

type FailedNode struct {
	SubID     uint16 `json:"sub_id"`
	UniqueKey uint64 `json:"unique_key"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Reason    string `json:"reason"`
}
