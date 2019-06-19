package couchdb

import (
	"fmt"
	"io"
)

type BulkCommitResponseRow struct {
	ID     string `json:"id"`
	Rev    string `json:"rev"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type BulkCommitResponse []BulkCommitResponseRow

type BulkCommit struct {
	AllOrNothing bool          `json:"all_or_nothing,omitempty"` // Not guaranteed on regular couchdb, not supported on cloudant. Generally avoid.
	NewEdits     *bool         `json:"new_edits,omitempty"`      // For replication
	Docs         []interface{} `json:"docs"`
}

// BulkUpdate accepts a commit with a list of updates to make to the DB, and returns a list of responses showing the status of each commit
func (db *CouchDB) BulkUpdate(c *BulkCommit) (*BulkCommitResponse, error) {
	for doc := range c.Docs {
		callWriteHook(doc)
	}
	var s BulkCommitResponse
	r, errCh := jsonifyDoc(c)
	req, err := db.createRequest("POST", "_bulk_docs", "", r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	code, cerr := couchDo(req, &s)
	if cerr != nil {
		return nil, cerr
	}
	if err = <-errCh; err != nil {
		return nil, err
	}
	if code != 201 {
		// FIXME Unexpected code. Do something?
	}
	return &s, nil
}

func (db *CouchDB) PutAttachment(docid string, docrev string, attachment io.Reader, attname string, ctype string) (*CouchSuccess, error) {
	var s CouchSuccess
	if docid == "" {
		return nil, MissingDocumentIDError
	}
	req, err := db.createRequest("PUT", escape_docid(docid)+"/"+escape_docid(attname), "rev="+docrev, attachment)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ctype)
	_, cerr := couchDo(req, &s)
	if cerr != nil {
		return nil, cerr
	}
	return &s, nil
}

func (db *CouchDB) PutDocument(doc interface{}, docid string) (*CouchSuccess, error) {
	return db.PutDocumentWithQueryString(doc, docid, "")
}

func (db *CouchDB) PutDocumentWithQueryString(doc interface{}, docid, queryString string) (*CouchSuccess, error) {
	var s CouchSuccess
	callWriteHook(doc)
	r, errCh := jsonifyDoc(doc)
	if docid == "" {
		return nil, MissingDocumentIDError
	}
	req, err := db.createRequest("PUT", escape_docid(docid), queryString, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	_, cerr := couchDo(req, &s)
	if cerr != nil {
		return nil, cerr
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	return &s, nil
}

func (db *CouchDB) PostDocument(doc interface{}) (*CouchSuccess, error) {
	return db.PostDocumentWithQueryString(doc, "")
}

func (db *CouchDB) PostDocumentWithQueryString(doc interface{}, queryString string) (*CouchSuccess, error) {
	var s CouchSuccess
	callWriteHook(doc)
	r, errCh := jsonifyDoc(doc)
	req, err := db.createRequest("POST", "", queryString, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	code, cerr := couchDo(req, &s)
	if cerr != nil {
		return nil, cerr
	}
	if err = <-errCh; err != nil {
		return nil, err
	}
	if code != 201 {
		// FIXME Unexpected code. Do something?
	}
	return &s, nil
}

func (db *CouchDB) DeleteDocument(docid, rev string) (*CouchSuccess, error) {
	var s CouchSuccess
	if docid == "" {
		return nil, MissingDocumentIDError
	}
	req, err := db.createRequest("DELETE", escape_docid(docid), fmt.Sprintf("rev=%s", rev), nil)
	if err != nil {
		return nil, err
	}
	code, cerr := couchDo(req, &s)
	if cerr != nil {
		return nil, cerr
	}
	if code != 200 {
		// FIXME Unexpected code. Do something?
	}
	return &s, nil
}

func callWriteHook(d interface{}) {
	if x, ok := d.(DocumentWithPreWriteHook); ok {
		x.CouchDocPreWrite()
	}
}
