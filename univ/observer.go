package univ

// Observer is an observer of body updates. See Body.AddObserver and Body.RemoveObserver
// for details on how to manage observers of a body.
type Observer interface {
	// BodyRotated is called on each observer whenever a body's rotation is changed.
	BodyRotated(body *Body)
	// BodyTranslated is called on each observer whenever a body's location is changed.
	BodyTranslated(body *Body)
}
