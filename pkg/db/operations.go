package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/dgraph-io/dgo/v2/y"
)

// This file should contain methods for interacting with the data easily.
// This includes reading and mutating data for the different data types

// GetScrape should recieve a dgraph client and a scrape struct,
// and return the official scrape struct from the database, complete with Uid for referencing
// if no such struct exists, then it returns an error
func GetScrape(c *dgo.Dgraph, scrape Scrape) (*Scrape, error) {
	if scrape.UID != "" {
		return getScrapeWithID(c, scrape)
	}
	return getScrapeWithoutID(c, scrape)
}

func getScrapeWithID(c *dgo.Dgraph, scrape Scrape) (*Scrape, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindScrape($uid: string) {
			findScrape(func: uid($uid)) {
				uid
				scrape.id
				scrape.last_scraped
				scrape.found_event {
					uid
					event.id
					event.title
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$uid"] = scrape.UID

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindScrape []Scrape `json:"findScrape"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindScrape) == 0 {
		return nil, fmt.Errorf("No Scrape found with id %s", scrape.UID)
	}

	return &r.FindScrape[0], nil
}

func getScrapeWithoutID(c *dgo.Dgraph, scrape Scrape) (*Scrape, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindScrapeNoID($id: int) {
			findScrapeNoID(func: eq(scrape.id, $id)) {
				uid
				scrape.id
				scrape.last_scraped
				scrape.found_event {
					uid
					event.id
					event.title
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = strconv.Itoa(scrape.ID)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindScrapeNoID []Scrape `json:"findScrapeNoID"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindScrapeNoID) == 0 {
		return nil, nil
	}

	return &r.FindScrapeNoID[0], nil
}

// UpsertScrape upserts the scrape struct into the database
func UpsertScrape(c *dgo.Dgraph, scrape Scrape) (*api.Response, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, err := json.Marshal(scrape)
	if err != nil {
		return nil, err
	}

	mu.SetJson = pb
	assigned, err := c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}
	return assigned, nil
}

// GetEvent should recieve a dgraph client and an event struct,
// and return the official event struct from the database, complete with Uid for referencing
// if no such event exists, then it returns an error
func GetEvent(c *dgo.Dgraph, event Event) (*Event, error) {
	if event.UID != "" {
		return getEventWithUID(c, event)
	}
	return getEventWithoutUID(c, event)
}

func getEventWithUID(c *dgo.Dgraph, event Event) (*Event, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindEvent($id: string) {
			findEvent(func: uid($id)) {
				uid
				event.id
				event.title
				event.description
				event.start_date
				event.end_date
				event.organiser {
					uid
					person.name
				}
				event.part_of_module {
					uid
					module.code
				}
				event.location {
					uid
					location.id
					location.name
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = event.UID

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindEvent []Event `json:"findEvent"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	if len(r.FindEvent) == 0 {
		return nil, fmt.Errorf("No Event found with uid %s", event.UID)
	}

	return &r.FindEvent[0], nil
}

func getEventWithoutUID(c *dgo.Dgraph, event Event) (*Event, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindEventNoUID($id: string) {
			findEvent(func: eq(event.id, $id)) {
				uid
				event.id
				event.title
				event.description
				event.start_date
				event.end_date
				event.organiser {
					uid
					person.name
				}
				event.part_of_module {
					uid
					module.code
				}
				event.location {
					uid
					location.id
					location.name
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = event.ID

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindEvent []Event `json:"findEvent"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindEvent) == 0 {
		return nil, nil
	}

	return &r.FindEvent[0], nil
}

// UpsertEvent upserts the event struct into the database
func UpsertEvent(c *dgo.Dgraph, event Event) (*api.Response, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, jsonErr := json.Marshal(event)
	if jsonErr != nil {
		return nil, jsonErr
	}

	mu.SetJson = pb
	assigned, upsertErr := c.NewTxn().Mutate(ctx, mu)
	if upsertErr != nil {
		if upsertErr == y.ErrAborted {
		}
		return nil, upsertErr
	}
	return assigned, nil
}

//GetLocationFromKentSlug returns a matching location from the slug kent uses internally
func GetLocationFromKentSlug(c *dgo.Dgraph, slug string) (*Location, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindLocationFromSlug($id: string) {
			findLocation(func: eq(location.id, $id)) {
				uid
				location.id
				location.name
				location.disabled_access
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = slug

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindLocation []Location `json:"findLocation"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindLocation) == 0 {
		return nil, nil
	}

	return &r.FindLocation[0], nil
}

// UpsertLocation upserts the location struct into the database
func UpsertLocation(c *dgo.Dgraph, loc Location) (*api.Response, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, err := json.Marshal(loc)
	if err != nil {
		return nil, err
	}

	mu.SetJson = pb
	assigned, err := c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}
	return assigned, nil
}

// CountNodesWithField returns the number of nodes which contain the location.id field
// this is a good indicator of the number of nodes of a certain type
// had to modify due to issues with variable passing
func CountNodesWithField(c *dgo.Dgraph, f string) (*int, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()

	q :=
		`query Count {
			nodeCount(func: has(location.id)) {
				nodeCount: count(uid)
			}
		}
		`

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	type Root struct {
		NodeCount []struct {
			NodeCount int `json:"nodeCount"`
		} `json:"nodeCount"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	return &r.NodeCount[0].NodeCount, nil
}