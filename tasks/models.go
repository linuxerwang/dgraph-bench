package tasks

const (
	TypePerson = iota
)

type Person struct {
	Uid       string `json:"uid,omitempty"`
	Name      string `json:"name,omitempty"`
	Xid       string `json:"xid,omitempty"`
	Type      int    `json:"type,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}
