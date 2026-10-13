package utils

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	GD       int    `json:"gd,omitempty"`
}

type ArgonUser struct {
	Account  int    `json:"account_id"`         // Player account ID
	User     int    `json:"user_id,omitempty"`  // Player user ID
	Username string `json:"username,omitempty"` // Player username
	Token    string `json:"authtoken"`          // Authorization token
}

type ArgonValidation struct {
	Valid     bool   `json:"valid"`                // If user is valid
	ValidWeak bool   `json:"valid_weak,omitempty"` // If user is just barely valid
	Username  string `json:"username,omitempty"`   // Player username
	Cause     string `json:"cause"`                // Cause for invalidation if any
}
