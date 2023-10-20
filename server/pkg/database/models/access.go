package models

const (
	RAcess   = "R"
	RWAccess = "RW"
)

type Access struct {
	Granter  string `json:"granter" prop:"granter"`
	Receiver string `json:"receiver" prop:"receiver"`
	Level    string `json:"level" prop:"level"`
}
