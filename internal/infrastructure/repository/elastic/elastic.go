package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"log/slog"
	"nearestPlaces/internal/entity"
	"nearestPlaces/internal/lib/logger/sl"
	"net/http"
)

type Storage struct {
	log    *slog.Logger
	client *elasticsearch.Client
	index  string
}

func New(log *slog.Logger, es *elasticsearch.Client, index string) *Storage {
	result := &Storage{
		log:    log,
		client: es,
		index:  index,
	}
	return result
}

func (e *Storage) expandMaxResultWindow() error {
	const op = "infrastructure.repository.elastic.expandMaxResultWindow"
	log := e.log.With(
		slog.String("op", op),
	)
	settings := map[string]interface{}{
		"index": map[string]interface{}{
			"max_result_window": 20000,
		},
	}

	body, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	req := esapi.IndicesPutSettingsRequest{
		Index: []string{e.index},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		log.Error("failed to update settings", sl.Err(err))
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error("status is not 200")
		return errors.New("status is not 200")
	}
	return nil
}

func (e *Storage) GetPlaces(limit, offset int) ([]*entity.Restaurant, int, error) {
	const op = "infrastructure.repository.elastic.GetPlaces"
	log := e.log.With(
		slog.String("op", op),
	)
	query := map[string]interface{}{
		"size": limit,
		"from": offset,
	}
	body, err := json.Marshal(query)
	if err != nil {
		log.Error("failed to marshal query", sl.Err(err))
		return nil, 0, err
	}

	req := esapi.SearchRequest{
		Index:          []string{e.index},
		Body:           bytes.NewReader(body),
		TrackTotalHits: true,
	}
	resp, err := req.Do(context.Background(), e.client)
	if err != nil {
		log.Error("failed to search in index", sl.Err(err))
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, 0, err
	}

	var respBody map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, 0, err
	}

	hits := respBody["hits"].(map[string]interface{})["hits"].([]interface{})
	rests := make([]*entity.Restaurant, 0, len(hits))
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		placeBytes, err := json.Marshal(source)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal place: %s", err.Error())
		}
		rest := &entity.Restaurant{}
		if err := json.Unmarshal(placeBytes, rest); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal place: %s", err.Error())
		}
		rests = append(rests, rest)
	}
	totalHits := int(respBody["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	return rests, totalHits, nil
}

func (e *Storage) CreateIndex(mappings []byte) error {
	const op = "infrastructure.repository.elastic.CreateIndex"
	log := e.log.With(
		slog.String("op", op),
	)
	resp, err := e.client.Indices.Exists([]string{e.index})
	if err != nil {
		log.Error("failed to check if index exists", sl.Err(err))
		return fmt.Errorf("error while checking if index exists: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		resp, err = e.client.Indices.Delete([]string{e.index})
		if err != nil || resp.IsError() {
			log.Error("failed to delete index")
			return fmt.Errorf("error while deleting index: %s", err)
		}
		resp.Body.Close()
	}
	resp, err = e.client.Indices.Create(e.index, e.client.Indices.Create.WithBody(bytes.NewBuffer(mappings)))
	if err != nil || resp.IsError() {
		log.Error("failed to create index")
		return fmt.Errorf("error while creating index: %v", err)
	}
	err = e.expandMaxResultWindow()
	resp.Body.Close()
	return err
}

func (e *Storage) SaveData(data []*entity.Restaurant) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  e.index,
		Client: e.client,
	})
	if err != nil {
		return fmt.Errorf("error creating bulk indexer: %w", err)
	}
	for _, d := range data {
		clause, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("error marshalling data: %w", err)
		}
		err = bi.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			Body:       bytes.NewReader(clause),
			DocumentID: d.ID,
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					e.log.Error("saveData: error while indexing data", sl.Err(err))
				}
			},
		})
		if err != nil {
			return fmt.Errorf("saveData: error saving data: %w", err)
		}
	}
	if err := bi.Close(context.Background()); err != nil {
		return fmt.Errorf("saveData: error closing bulk indexer: %w", err)
	}
	return nil
}

func (e *Storage) GetClosest(lat, lon float64) ([]*entity.Restaurant, error) {
	const op = "infrastructure.repository.elastic.GetClosest"
	log := e.log.With(
		slog.String("op", op),
	)
	query := map[string]interface{}{
		"size": 3,
		"sort": map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": lat,
					"lon": lon,
				},
				"order":           "asc",
				"unit":            "km",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		},
	}
	body, err := json.Marshal(query)
	if err != nil {
		log.Error("failed to marshal query", sl.Err(err))
		return nil, err
	}
	req := esapi.SearchRequest{
		Index: []string{e.index},
		Body:  bytes.NewReader(body),
	}
	resp, err := req.Do(context.Background(), e.client)
	if err != nil {
		log.Error("failed to search in index", sl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		log.Error("failed to search in index")
		return nil, errors.New("error while search in index")
	}

	var respBody map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Error("failed to unmarshal response", sl.Err(err))
		return nil, err
	}

	hits := respBody["hits"].(map[string]interface{})["hits"].([]interface{})
	rests := make([]*entity.Restaurant, 0, len(hits))
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		placeBytes, err := json.Marshal(source)
		if err != nil {
			log.Error("failed to marshal place", sl.Err(err))
			return nil, err
		}
		rest := &entity.Restaurant{}
		if err := json.Unmarshal(placeBytes, rest); err != nil {
			log.Error("failed to unmarshal place", sl.Err(err))
			return nil, err
		}
		rests = append(rests, rest)
	}
	return rests, nil
}
