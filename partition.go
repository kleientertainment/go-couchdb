package couchdb

import "fmt"

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
// Encodes a AllDocsParams struct into a query string
func (p *AllDocsParams) Encode() (string, error) {
	return urlEncodeObject(*p)
}

type Partition struct {
	db *CouchDB
	partition string
}

func (db *CouchDB) Partition(partition string) *Partition {
	return &Partition{db: db, partition: partition}
}

func (p *Partition) buildPath(path string) string {
	return fmt.Sprintf("/_partition/%s/%s", p.partition, path)
}

// This is basically extracted from the viewHelper. Keeping the old code as it was to avoid breaking anything.
func (p *Partition) AllDocs(params AllDocsParams) (*AllDocsResult, error) {
	argstring, err := params.Encode()
	if err != nil {
		return nil, err
	}
	req, err := p.db.createRequest("GET", p.buildPath("_all_docs"), argstring, nil)
	if err != nil {
		return nil, err
	}

	results := new(AllDocsResult)
	if _, err = couchDo(req, results); err != nil {
		return nil, err
	}

	return results, nil
}
