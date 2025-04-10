package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"nearestPlaces/internal/entity"
	"os"
	"strconv"
)

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) ParseCSV(filename string) ([]*entity.Restaurant, error) {
	var res []*entity.Restaurant
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	r := csv.NewReader(bufio.NewReader(file))
	r.Comma = '\t'
	_, _ = r.Read()
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read line: %w", err)
		}
		data, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line \"%s\": %w", line, err)
		}
		res = append(res, data)
	}
	return res, nil
}

func parseLine(line []string) (*entity.Restaurant, error) {
	if len(line) != 6 {
		return nil, fmt.Errorf("failed to parse line \"%s\": wrong number of fields", line)
	}
	lon, err := strconv.ParseFloat(line[4], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Longitude: %s", line[4])
	}
	lat, err := strconv.ParseFloat(line[5], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Latitude: %s", line[5])
	}
	return &entity.Restaurant{
		ID:      line[0],
		Name:    line[1],
		Address: line[2],
		Phone:   line[3],
		Location: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{Lon: lon, Lat: lat}}, nil
}
