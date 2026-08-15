package match

// FormatDescriptor is match metadata owned by whoever created the match. The
// simulator does not interpret it and no game rule depends on it: the values
// are stored on the match and reported back in match summaries and the duel
// result webhook, so the caller can identify a match's format even after the
// caller itself has restarted.
type FormatDescriptor struct {
	// ID identifies the format so the caller can look its rules up.
	ID string
	// Name is a display label captured when the match was created, so it
	// survives the format being renamed or deleted.
	Name string
}
