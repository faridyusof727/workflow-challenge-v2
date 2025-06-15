package edge

import "time"

type Edge struct {
	ID        string         `json:"id"`
	Source    string         `json:"source"`
	Target    string         `json:"target"`
	Kind      string         `json:"type"`
	Animated  bool           `json:"animated"`
	Style     map[string]any `json:"style"`
	Label     string         `json:"label"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
