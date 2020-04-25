package match

// Context is events passed down to cards, allowing them to perform actions
// without having a direct reference to the match, players etc
type Context struct {
	Match  Match
	Event  string
	cancel bool
}

// HandlerFunc is a function with a match context as argument
type HandlerFunc func(c *Context)

// ScheduleAfter allows you to run the logic at the end of the context flow,
// after the default behaviour
func (c *Context) ScheduleAfter() {

}

// PreventDefault stops the context flow, cancelling the default behaviour
func (c *Context) PreventDefault() {

}
