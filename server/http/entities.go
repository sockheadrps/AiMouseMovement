package http

type mousePoint struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Time float32 `json:"time"`
}

type dataSet struct {
	WindowHeight int32        `json:"window-height"`
	WindowWidth  int32        `json:"window-width"`
	MouseArray   []mousePoint `json:"mouse-array"`
}

type validator struct {
	User string `json:"username"`
	Pwd  string `json:"password"`
}

type clientUuid struct {
	Uuid string `json:"uuid"`
}

type aprovedDataSet struct {
	Id           string       `json:"_id"`
	Uuid         string       `json:"uuid"`
	WindowHeight int32        `json:"window-height"`
	WindowWidth  int32        `json:"window-width"`
	MouseArray   []mousePoint `json:"mouse-array"`
}

type removeDataSet struct {
	Id   string `json:"_id"`
	Uuid string `json:"uuid"`
}

type RandomDocument struct {
	Id           string       `json:"_id" bson:"_id"`
	WindowHeight int32        `json:"window-height" bson:"windowheight"`
	WindowWidth  int32        `json:"window-width" bson:"windowwidth"`
	MouseArray   []mousePoint `json:"mouse-array" bson:"mousearray"`
	Uuid         string       `json:"uuid"`
}
