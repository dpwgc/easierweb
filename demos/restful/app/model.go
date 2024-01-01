package app

// MemberCommand add or edit request body
type MemberCommand struct {
	Name   string  `json:"name"`
	Mobile string  `json:"mobile"`
	Weight float32 `json:"weight"`
}

// MemberQuery list query request params
type MemberQuery struct {
	Name   string `mapstructure:"name"`
	Mobile string `mapstructure:"mobile"`
}

// MemberDTO get or list response object
type MemberDTO struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Mobile string  `json:"mobile"`
	Weight float32 `json:"weight"`
}
