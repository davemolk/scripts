package main

type queryData struct {
	base        string
	controlName string
	query       string
	queryString string
	spacer      string
	terms       []string
}

type queryMap struct {
	queries map[string]queryData
}

func (f *fof) makeQueryData() []*queryData {
	var qd []*queryData

	bing := &queryData{
		base:        "https://bing.com/search?q=",
		controlName: "- Search",
		query:       f.config.query,
		spacer:      "%20",
		terms:       f.terms,
	}
	qd = append(qd, bing)

	// blocks a lot
	google := &queryData{
		base:        "https://www.google.com/search?q=",
		controlName: "- Google Search",
		query:       f.config.query,
		spacer:      "+",
		terms:       f.terms,
	}

	_ = google

	yahoo := &queryData{
		base:        "https://search.yahoo.com/search?p=",
		controlName: "- Yahoo Search Results",
		query:       f.config.query,
		spacer:      "+",
		terms:       f.terms,
	}

	qd = append(qd, yahoo)

	return qd
}
