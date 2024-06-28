package ctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wyubin/llm-svc/utils/str"
)

const (
	pathListModel = "/api/tags"
	pathEmbedding = "/api/embeddings"
	pathChat      = "/api/chat"
	nameChatTmpl  = "rag"
)

type Gpt struct {
	uri     string
	dirTMPL string
}

func (s *Gpt) SetDirTmpl(dirTmpl string) {
	s.dirTMPL = dirTmpl
}

// check model exist, if not, pull model with name
func (s *Gpt) HasModel(name string) (bool, error) {
	// get list of models
	models, err := s.GetModels()
	if err != nil {
		return false, err
	}
	for _, model := range models {
		if model.Name == name || strings.Split(model.Model, ":")[0] == name {
			return true, nil
		}
	}
	return false, nil
}

func (s *Gpt) GetModels() ([]RespModel, error) {
	// get list of models
	resp, err := http.Get(fmt.Sprintf("%s%s", s.uri, pathListModel))
	if err != nil {
		return nil, err
	}
	models := RespModels{}
	err = json.NewDecoder(resp.Body).Decode(&models)
	if err != nil {
		return nil, err
	}
	return models.Models, nil
}

func (s *Gpt) Embedding(text, nameModel string) ([]float32, error) {
	data := map[string]string{
		"model":  nameModel,
		"prompt": text,
	}
	postBody, _ := json.Marshal(data)
	resp, err := http.Post(fmt.Sprintf("%s%s", s.uri, pathEmbedding), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}
	embedding := RespEmbedding{}
	err = json.NewDecoder(resp.Body).Decode(&embedding)
	if err != nil {
		return nil, err
	}
	return embedding.Embedding, nil
}

// chat llm with contexts
func (s *Gpt) ChatWithContexts(nameModel, question string, contexts []string) (string, error) {
	pathTmpl := filepath.Join(s.dirTMPL, nameChatTmpl+".tmpl")
	fileTmpl, err := os.Open(pathTmpl)
	if err != nil {
		return "", err
	}
	bufTmpl := &strings.Builder{}
	io.Copy(bufTmpl, fileTmpl)
	kbSystem, err := str.Tmpl2Str(bufTmpl.String(), contexts)
	if err != nil {
		return "", err
	}
	// prepare messages
	msgs := []map[string]string{
		{"role": "system", "content": kbSystem},
		{"role": "user", "content": "Question: " + question},
	}
	postBody, _ := json.Marshal(map[string]interface{}{
		"model":    nameModel,
		"messages": msgs,
		"stream":   false,
	})
	resp, err := http.Post(fmt.Sprintf("%s%s", s.uri, pathChat), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return "", err
	}
	data := RespChat{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Message.Content, nil
}

func NewGpt(uri string) (*Gpt, error) {
	gpt := Gpt{
		uri: uri,
	}
	// check service available
	_, err := gpt.GetModels()
	if err != nil {
		return nil, fmt.Errorf("service %s not available", uri)
	}
	return &gpt, nil
}
