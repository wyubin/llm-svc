package ctl

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conf struct {
	path string
}

func (s *Conf) LoadConfig() config {
	file, _ := os.Open(s.path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	res := config{}
	decoder.Decode(&res)
	return res
}

func (s *Conf) SaveConfig(config config) error {
	file, _ := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func (s *Conf) GetCollInfo(name string) (Collection, error) {
	config := s.LoadConfig()
	for _, coll := range config.Collections {
		if coll.Name == name {
			return coll, nil
		}
	}
	return Collection{}, fmt.Errorf("no such collection: %s", name)
}

func (s *Conf) UpdateColl(coll Collection) error {
	config := s.LoadConfig()
	needAppend := true
	for i, c := range config.Collections {
		if c.Name == coll.Name {
			config.Collections[i] = coll
			needAppend = false
			break
		}
	}
	if needAppend {
		config.Collections = append(config.Collections, coll)
	}
	return s.SaveConfig(config)
}

func (s *Conf) UpdateModel(model Model) error {
	config := s.LoadConfig()
	needAppend := true
	for i, m := range config.Models {
		if m.Name == model.Name {
			config.Models[i] = model
			needAppend = false
			break
		}
	}
	if needAppend {
		config.Models = append(config.Models, model)
	}
	return s.SaveConfig(config)
}

type config struct {
	Models      []Model      `json:"models"`
	Collections []Collection `json:"collections"`
}

type Model struct {
	Name      string `json:"name"`
	VectorLen int    `json:"vector-length"`
}

type Collection struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

func NewConf(path string) (*Conf, error) {
	// check path exists, if not, create and save empty
	_, err := os.Stat(path)
	if err != nil {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.Encode(config{})
	}
	return &Conf{path: path}, nil
}
