package infra

import (
	"encoding/json"
	"io"
	"os"
)

type Store[T any] struct {
	filepath string
}

func NewStore[T any](filepath string) *Store[T] {
	return &Store[T]{
		filepath: filepath,
	}
}

func (s *Store[T]) Read() ([]T, error) {
	f, err := os.Open(s.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var result []T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Store[T]) Write(data []T) error {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	f, err := os.Create(s.filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
