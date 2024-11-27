package mappers

type RawTLE struct {
	NoradID string `json:"norad_id"`
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
}
