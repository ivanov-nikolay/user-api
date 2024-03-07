package filters

type Filter struct {
	Gender           string
	Status           string
	FullName         string
	AttributesToSort string
	SortAsk          bool
	SortDesc         bool
	Limit            int
	Offset           int
}
