package types

type TogglResponse struct {
	Groups []TogglResponseGroup `json:"groups,omitempty"`
}

type TogglResponseGroup struct {
	SubGroup []TogglResponseSubGroup `json:"sub_groups,omitempty"`
}

type TogglResponseSubGroup struct {
	Title   string `json:"title"`
	Seconds int    `json:"seconds"`
}
