package store

type stringError struct {
	payload string
}

func (err stringError) Error() string {
	return "store: " + err.payload
}
