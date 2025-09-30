package model

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	OwnerID     string `json:"owner_id"`
}

func (e *Event) Validate() bool {
	if e.ID == "" || e.Title == "" || e.Date == "" || e.OwnerID == "" {
		return false
	}
	return true
}
