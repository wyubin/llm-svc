package ctl

import (
	"fmt"

	repoDB "github.com/wyubin/llm-svc/pkg/repo/vector"
	interDB "github.com/wyubin/llm-svc/utils/repo/vector"
)

const (
	KEY_EMBEDDING = "context"
)

type DB struct {
	conn interDB.Conn
	coll interDB.Coll
}

func (s *DB) HasColl(name string) bool {
	_, err := s.conn.GetColl(name)
	return err == nil
}

func (s *DB) CreateColl(name string, vectorSize int) error {
	_, err := s.conn.NewColl(name, uint64(vectorSize))
	return err
}

func (s *DB) ChooseColl(name string) error {
	coll, err := s.conn.GetColl(name)
	if err != nil {
		return err
	}
	s.coll = coll
	return nil
}

// insert one record
func (s *DB) InsertPoint(embedding interface{}, meta map[string]string) error {
	if s.coll == nil {
		return fmt.Errorf("please choose a collection first")
	}
	_, err := s.coll.InsertOne(embedding, meta)
	return err
}

// vector search
func (s *DB) VectorSearch(embedding interface{}, limit int) ([]string, error) {
	if s.coll == nil {
		return nil, fmt.Errorf("please choose a collection first")
	}
	points := []interDB.Point{}
	err := s.coll.VectorSearch(embedding, &points, map[string]interface{}{"limit": limit})
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, p := range points {
		context, found := p.GetMeta()[KEY_EMBEDDING]
		if !found {
			continue
		}
		res = append(res, context)
	}
	return res, nil
}

func (s *DB) Close() error {
	return s.conn.Close()
}

func NewDB(uri string) (*DB, error) {
	conn, err := repoDB.NewConn(uri)
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}
