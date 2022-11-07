package fine

// This function recieves the application it has been called from.
type CronFunc func(app *App)

type QueuedFunc struct {
	Func         CronFunc // The function to be called.
	Left         float64  // Seconds left until this function is called.
	shouldRepeat bool
	startTime    float64
}

// Runs a function after a specified amount of seconds.
// This is not precise, the functions are updated every frame.
// This will only start counting after app.Run has been called.
// The function recieves the application it has been called from.
func (app *App) After(seconds float64, afterFunction CronFunc) {
	app.QueuedFunctions = append(app.QueuedFunctions, &QueuedFunc{
		Func: afterFunction,
		Left: seconds,
	})
}

// Starts running a function every n amount of seconds.
func (app *App) Every(seconds float64, everyFunc CronFunc) {
	app.QueuedFunctions = append(app.QueuedFunctions, &QueuedFunc{
		Func:         everyFunc,
		Left:         seconds,
		shouldRepeat: true,
		startTime:    seconds,
	})
}
