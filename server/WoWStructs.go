package server

type Hero struct {
	Name   string `json:"name"`
	Race   bool   `json:"race"`
	Gender bool   `json:"gender"`
	Level  int    `json:"level"`
	Class  int    `json:"class"`
}
type HeroPosition struct {
	PositionX float32 `json:"position_x"`
	PositionY float32 `json:"position_y"`
	PositionZ float32 `json:"position_z"`
}
type Home struct {
	Alliance HeroPosition `json:"alliance"`
	Horde    HeroPosition `json:"horde"`
}
