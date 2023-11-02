package http

type mousePoint struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Time int64   `json:"time"`
}

type dataSet struct {
	WindowHeight int32        `json:"window-height"`
	WindowWidth  int32        `json:"window-width"`
	MouseArray   []mousePoint `json:"mouse-array"`
}