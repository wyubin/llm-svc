package vector

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/google/uuid"
	"github.com/wyubin/llm-svc/utils/maptool"
	"github.com/wyubin/llm-svc/utils/repo/vector"

	pb "github.com/qdrant/go-client/qdrant"
)

var (
	waitUpsert         = true
	limitSearch uint64 = 10
)

type QdrantColl struct {
	client   pb.PointsClient
	nameColl string
}

func (s *QdrantColl) VectorSearch(embeddings interface{}, records interface{}, args ...map[string]interface{}) error {
	mergeArgs := map[string]interface{}{}
	maptool.Update(mergeArgs, args...)
	vectorFloats, err := inter2floats(embeddings)
	if err != nil {
		return err
	}
	searchRule := pb.SearchPoints{
		CollectionName: s.nameColl,
		Vector:         vectorFloats,
		Limit:          limitSearch,
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
	}
	// limit
	argLimit, found := mergeArgs["limit"]
	if found {
		searchRule.Limit = uint64(argLimit.(int))
	}
	res, err := s.client.Search(context.Background(), &searchRule)
	if err != nil {
		return err
	}
	points := []vector.Point{}
	for _, hit := range res.GetResult() {
		point := QdrantPoint{scoredPoint: hit}
		points = append(points, &point)
	}

	reflect.ValueOf(records).Elem().Set(reflect.ValueOf(points))
	return nil
}

func (s *QdrantColl) InsertOne(vector interface{}, meta map[string]string) (string, error) {
	// use uuid as id
	uuidV4, _ := uuid.NewRandom()
	idPd := pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{Uuid: uuidV4.String()},
	}
	vectorFloats, err := inter2floats(vector)
	if err != nil {
		return "", err
	}
	vectorPb := pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: vectorFloats}}}
	metaPb := map[string]*pb.Value{}
	for k, v := range meta {
		metaPb[k] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: v}}
	}
	_, err = s.client.Upsert(context.Background(), &pb.UpsertPoints{
		CollectionName: s.nameColl,
		Wait:           &waitUpsert,
		Points: []*pb.PointStruct{
			{
				Id:      &idPd,
				Vectors: &vectorPb,
				Payload: metaPb,
			},
		},
	})
	if err != nil {
		return "", err
	}
	// use uuid as id
	return uuidV4.String(), nil
}

func (s *QdrantColl) DeleteByID(id string) error {
	idPd := pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{Uuid: id},
	}
	_, err := s.client.Delete(context.Background(), &pb.DeletePoints{
		CollectionName: s.nameColl,
		Points: &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Points{
				Points: &pb.PointsIdsList{
					Ids: []*pb.PointId{&idPd},
				},
			},
		},
	})
	return err
}

func (s *QdrantColl) CreateIndexes(nameCols ...string) error {
	return nil
}

func inter2floats(inter interface{}) ([]float32, error) {
	vectorStr, _ := json.Marshal(inter)
	vectorFloats := []float32{}
	err := json.Unmarshal([]byte(vectorStr), &vectorFloats)
	if err != nil {
		return nil, err
	}
	return vectorFloats, nil
}
