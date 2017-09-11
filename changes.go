package couchdb

import (
	"encoding/json"
	"fmt"
)

type DocRev struct {
	ID      string          `json:"id"`
	Seq     interface{}     `json:"seq"`
	Doc     json.RawMessage `json:"doc"`
	Deleted bool            `json:"deleted"`
	Changes []Rev           `json:"changes"`
}

type Rev struct {
	Rev string `json:"rev"`
}

type ChangesArgs struct {
	Since       interface{}     `urlencode:"since"`
	Limit       uint64          `urlencode:"limit"`
	Descending  bool            `urlencode:"descending"`
	Feed        UnescapedString `urlencode:"feed"`
	Heartbeat   uint64          `urlencode:"heartbeat"`
	Timeout     uint64          `urlencode:"timeout"`
	Filter      string          `urlencode:"filter"`
	IncludeDocs bool            `urlencode:"include_docs"`
	Style       string          `urlencode:"style"`
	SeqInterval uint64          `urlencode:"seq_interval"`
}

func (a *ChangesArgs) Encode() (string, error) {
	return urlEncodeObject(*a)
}

type Changes struct {
	Results []DocRev    `json:"results"`
	LastSeq interface{} `json:"last_seq"`
}

func (db *CouchDB) Changes(args ChangesArgs) (*Changes, error) {
	if args.Feed == "continuous" {
		return nil, fmt.Errorf("Changes is for non-continuous feeds. Try ContinuousChanges instead")
	}
	panic("Not implemented yet")
}

// ContinuousChanges starts a feed=continuous view of the _changes feed for the DB.
// Each change will be emitted from the *DocRev channel until the server hangs
// up, at which time the DocRev channel will be closed and the error channel
// will spit out the appropriate error
func (db *CouchDB) ContinuousChanges(args ChangesArgs) (<-chan *DocRev, <-chan error) {
	c := make(chan *DocRev)
	e := make(chan error, 1)
	args.Feed = "continuous"
	argsstring, err := args.Encode()
	if err != nil {
		e <- err
		return nil, e
	}
	req, err := db.createRequest("GET", "_changes", argsstring, nil)
	if err != nil {
		e <- err
		return nil, e
	}
	r, err := DefaultClient.Do(req)
	if err != nil {
		e <- err
		return nil, e
	}
	if r.StatusCode != 200 {
		r.Body.Close()
		e <- responseToCouchError(r)
		return nil, e
	}
	j := json.NewDecoder(r.Body)
	go func() {
		defer close(c)
		defer r.Body.Close()
		for {
			var r DocRev
			if err := j.Decode(&r); err != nil {
				e <- err
				return
			}
			if r.Seq == 0 {
				e <- fmt.Errorf("Sequence number was not set, or set to 0", r.Seq)
				return
			}
			c <- &r
		}
	}()
	return c, e
}
