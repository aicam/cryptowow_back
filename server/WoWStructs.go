package server

type Hero struct {
	Name   string `json:"name"`
	Race   bool   `json:"race"`
	Gender bool   `json:"gender"`
	Level  int    `json:"level"`
	Class  int    `json:"class"`
}

type Home struct {
}
