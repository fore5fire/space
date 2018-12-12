package draw

type Drawable interface {
	Draw(*GLState)
}

type GLState struct {
}
