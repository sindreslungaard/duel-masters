package match

type Format string

const (
	RegularFormat Format = "regular"
	RandomFormat  Format = "random"
)

func FormatFromStr(input string) Format {
	switch input {
	case "regular":
		return RegularFormat
	case "random":
		return RandomFormat
	default:
		return RegularFormat
	}
}
