package search

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearch struct {
	Client *elasticsearch.Client
	index  string
	Alias  string
}

func New(addresses []string) (*ElasticSearch, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticSearch{
		Client: client,
	}, nil
}

func (e *ElasticSearch) CreateIndex(index string) error {
	e.index = index
	e.Alias = index + "_alias"

	res, err := e.Client.Indices.Exists([]string{e.index})
	if err != nil {
		return fmt.Errorf("cannot check index existence: %w", err)
	}
	if res.StatusCode == 200 {
		return nil
	}
	if res.StatusCode != 404 {
		return fmt.Errorf("error in index existence response: %s", res.String())
	}

	res, err = e.Client.Indices.Create(e.index)
	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error in index creation response: %s", res.String())
	}

	res, err = e.Client.Indices.PutAlias([]string{e.index}, e.Alias)
	if err != nil {
		return fmt.Errorf("cannot create index alias: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error in index alias creation response: %s", res.String())
	}

	return nil
}
