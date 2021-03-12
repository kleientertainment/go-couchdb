package couchdb

type AllDocsParams struct {
	Key            interface{}     `urlencode:"key"`
	Conflicts 	   bool 		   `urlencode:"conflicts"`
	Descending     bool            `urlencode:"descending"`
	EndKey         interface{}     `urlencode:"endkey"`
	EndKey_DocID   string          `urlencode:"endkey_docid"`
	IncludeDocs    bool            `urlencode:"include_docs"`
	InclusiveEnd   *bool           `urlencode:"inclusive_end"` // Because the default is true
	Limit          uint            `urlencode:"limit"`
	Skip           uint            `urlencode:"skip"`
	StartKey       interface{}     `urlencode:"startkey"`
	StartKey_DocID string          `urlencode:"startkey_docid"`
	UpdateSeq      bool            `urlencode:"update_seq"`
}

type Partition struct {
	db *CouchDB
	partition string
}

func (db *CouchDB) Partition(partition string) *Partition {
	return &Partition{db: db, partition: partition}
}

func (p *Partition) AllDocs(params AllDocsParams) (*AllDocsResult, error) {

	return nil, nil
}