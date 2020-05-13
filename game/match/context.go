package match

// Context is events passed down to cards, allowing them to perform actions
// without having a direct reference to the match, players etc
type Context struct {
	Match   *Match
	Event   interface{}
	cancel  bool
	postFxs []func()
}

// HandlerFunc is a function with a match context as argument
type HandlerFunc func(card *Card, c *Context)

// NewContext returns a new match context
func NewContext(m *Match, e interface{}) *Context {

	ctx := &Context{
		Match:  m,
		Event:  e,
		cancel: false,
	}

	return ctx

}

// ScheduleAfter allows you to run the logic at the end of the context flow,
// after the default behaviour
func (c *Context) ScheduleAfter(handlers ...func()) {
	c.postFxs = append(c.postFxs, handlers...)
}

// InterruptFlow stops the context flow, cancelling the default behaviour
func (c *Context) InterruptFlow() {
	c.cancel = true
}
