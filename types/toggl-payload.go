package types

type TogglPayload struct {
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	IncludeTimeEntry bool   `json:"include_time_entry_ids"`
	HideRates        bool   `json:"hide_rates"`
	HideAmounts      bool   `json:"hide_amounts"`
}
