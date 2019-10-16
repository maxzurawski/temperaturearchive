package handlers

type Source int

const (
	Temperature Source = iota
)

var names = []string{
	"temperature",
}

func (s Source) String() string {
	return names[s]
}
