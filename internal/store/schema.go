package store

type Command struct {
	Text     string
	UseCount int
	LastUsed int64
}

type Exclusion struct {
	Pattern   string
	CreatedAt int64
}
