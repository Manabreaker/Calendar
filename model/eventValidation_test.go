package model

import "testing"

func TestEventValidate(t *testing.T) {
	tests := []struct {
		name  string
		event Event
		want  bool
	}{
		{
			name: "valid event",
			event: Event{
				ID:          "1",
				Title:       "Meeting",
				Description: "Project meeting",
				Date:        "2024-06-01",
				OwnerID:     "user1",
			},
			want: true,
		},
		{
			name: "missing title",
			event: Event{
				ID:          "2",
				Title:       "",
				Description: "No title event",
				Date:        "2024-06-02",
				OwnerID:     "user2",
			},
			want: false,
		},
		{
			name: "missing date",
			event: Event{
				ID:          "3",
				Title:       "No Date",
				Description: "Event without date",
				Date:        "",
				OwnerID:     "user3",
			},
			want: false,
		},
		{
			name: "missing owner ID",
			event: Event{
				ID:          "4",
				Title:       "No Owner",
				Description: "Event without owner",
				Date:        "2024-06-03",
				OwnerID:     "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.Validate(); got != tt.want {
				t.Errorf("Event.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
