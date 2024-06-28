package vector

import (
	"reflect"

	pb "github.com/qdrant/go-client/qdrant"
)

type QdrantPoint struct {
	scoredPoint *pb.ScoredPoint
}

func (s *QdrantPoint) GetID() string {
	return s.scoredPoint.GetId().GetUuid()
}

func (s *QdrantPoint) GetVector(record interface{}) error {
	data := s.scoredPoint.Vectors.GetVector().GetData()
	resultVar := reflect.ValueOf(data)
	reflect.ValueOf(record).Elem().Set(resultVar)
	return nil
}

func (s *QdrantPoint) GetMeta() map[string]string {
	meta := map[string]string{}
	for k, v := range s.scoredPoint.GetPayload() {
		meta[k] = v.GetStringValue()
	}
	return meta
}

func (s *QdrantPoint) GetScore() float32 {
	return s.scoredPoint.GetScore()
}
