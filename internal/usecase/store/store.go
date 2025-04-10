package store

import (
	"fmt"
	"log/slog"
	"nearestPlaces/internal/entity"
	"nearestPlaces/internal/lib/config"
	"nearestPlaces/internal/lib/logger/sl"
)

type SchemaReader interface {
	ReadMappings(filename string) ([]byte, error)
}

type Storage interface {
	CreateIndex(mappings []byte) error
	SaveData(data []*entity.Restaurant) error
}

type CSVParser interface {
	ParseCSV(filename string) ([]*entity.Restaurant, error)
}

type UseCase struct {
	log          *slog.Logger
	cfg          *config.Config
	schemaReader SchemaReader
	csvParser    CSVParser
	storage      Storage
}

func New(log *slog.Logger, cfg *config.Config, reader SchemaReader, parser CSVParser, storage Storage) *UseCase {
	return &UseCase{
		log:          log,
		cfg:          cfg,
		schemaReader: reader,
		csvParser:    parser,
		storage:      storage,
	}
}

func (u *UseCase) CreateIndexWithMapping() error {
	const op = "usecase.store.createIndexWithMapping"
	log := u.log.With(
		slog.String("op", op),
	)
	mappings, err := u.schemaReader.ReadMappings(u.cfg.SchemaPath)
	if err != nil {
		log.Error("failed to read mappings: ", sl.Err(err))
		return err
	}
	logMsg := fmt.Sprintf("successfully read mappings from %s", u.cfg.SchemaPath)
	log.Info(logMsg)

	err = u.storage.CreateIndex(mappings)
	if err != nil {
		log.Error("failed to create index: ", sl.Err(err))
		return err
	}
	log.Info("successfully created index")
	return nil
}

func (u *UseCase) UploadPlaces() error {
	const op = "usecase.store.fillIndex"
	log := u.log.With(
		slog.String("op", op),
	)
	data, err := u.csvParser.ParseCSV(u.cfg.DataPath)
	if err != nil {
		log.Error("failed to parse data: ", sl.Err(err))
		return err
	}
	log.Info("data parsed successfully")

	err = u.storage.SaveData(data)
	if err != nil {
		log.Error("failed to save data: ", sl.Err(err))
	}
	log.Info("successfully saved data")
	return nil
}
