package fine

// This function recieves the application it has been called from.
type AfterFunc func(*App)

type QueuedFunc struct {
	Func AfterFunc // The function to be called.
	Left float64   // Seconds left until this function is called.
}

// Runs a function after a specified amount of seconds.
// This is not precise, the functions are updated every frame.
// This will only start counting after app.Run has been called.
// The function recieves the application it has been called from.
func (app *App) After(seconds float64, afterFunction AfterFunc) {
	app.QueuedFunctions = append(app.QueuedFunctions, &QueuedFunc{
		Func: afterFunction,
		Left: seconds,
	})
}
