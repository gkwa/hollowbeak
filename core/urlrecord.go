package core

type urlRecord struct {
	URL string
}

func newURLRecord(rawURL string) urlRecord {
	return urlRecord{URL: rawURL}
}
