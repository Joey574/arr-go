package qbit

type Content struct {
	Index        int     `json:"index"`
	Name         string  `json:"name"`
	Size         int     `json:"size"`
	Progress     float64 `json:"progress"`
	Priority     int     `json:"priority"`
	IsSeed       bool    `json:"is_seed"`
	PieceRange   []int   `json:"piece_range"`
	Availability float64 `json:"availability"`
}
