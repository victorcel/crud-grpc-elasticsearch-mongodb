package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/models"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/search"
	"io"
	"time"
)

type UserStorage struct {
	elastic search.ElasticSearch
	timeout time.Duration
}

func NewUserStorage(elastic search.ElasticSearch) (UserStorage, error) {
	return UserStorage{
		elastic: elastic,
		timeout: time.Second * 10,
	}, nil
}

func (c *UserStorage) InsertUserElastic(ctx context.Context, UserElastic *models.UserElastic) (string, error) {
	bdy, err := json.Marshal(UserElastic)

	if err != nil {
		return "", fmt.Errorf("insert: marshall: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      c.elastic.Alias,
		DocumentID: UserElastic.ID,
		Body:       bytes.NewReader(bdy),
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := req.Do(ctx, c.elastic.Client)

	if err != nil {
		return "", fmt.Errorf("insert: request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.StatusCode == 409 {
		return "", errors.New("conflict")
	}

	if res.IsError() {
		return "", fmt.Errorf("insert: response: %s", res.String())
	}

	var (
		response interface{}
		body     document
	)
	body.Id = &response

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("find one: decode: %w", err)
	}

	return fmt.Sprint(response), nil
}

func (c *UserStorage) GetUserElasticByID(ctx context.Context, id string) (*models.UserElastic, error) {
	req := esapi.GetRequest{
		Index:      c.elastic.Alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := req.Do(ctx, c.elastic.Client)
	if err != nil {
		return &models.UserElastic{}, fmt.Errorf("find one: request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.StatusCode == 404 {
		return &models.UserElastic{}, errors.New("not found")
	}

	if res.IsError() {
		return &models.UserElastic{}, fmt.Errorf("find one: response: %s", res.String())
	}

	var (
		UserElastic models.UserElastic
		body        document
	)
	body.Source = &UserElastic

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return &models.UserElastic{}, fmt.Errorf("find one: decode: %w", err)
	}

	return &UserElastic, nil
}

func (c *UserStorage) UpdateUserElastic(ctx context.Context, UserElastic models.UserElastic) error {
	bdy, err := json.Marshal(UserElastic)
	if err != nil {
		return fmt.Errorf("update: marshall: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      c.elastic.Alias,
		DocumentID: UserElastic.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bdy))),
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := req.Do(ctx, c.elastic.Client)
	if err != nil {
		return fmt.Errorf("update: request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.StatusCode == 404 {
		return errors.New("not found")
	}

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}

func (c *UserStorage) DeleteUserElastic(ctx context.Context, id string) error {

	req := esapi.DeleteRequest{
		Index:      c.elastic.Alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := req.Do(ctx, c.elastic.Client)
	if err != nil {
		return fmt.Errorf("delete: request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.StatusCode == 404 {
		return errors.New("not found")
	}

	if res.IsError() {
		return fmt.Errorf("delete: response: %s", res.String())
	}

	return nil
}

type document struct {
	Source interface{} `json:"_source"`
	Id     interface{} `json:"_id"`
}
