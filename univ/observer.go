package univ

// Observer is an observer of body updates. See Body.AddObserver and Body.RemoveObserver
// for details on how to manage observers of a body.
type Observer interface {
	// BodyUpdated will be called periodically for each body being observed, and should
	// be used to perform any needed processing.
	BodyUpdated(body *Body, secondsElapsed float32)
}
