package vector

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/qdrant/go-client/qdrant"
	"github.com/wyubin/go-utils/repo/vector"
)

var (
	defaultSegmentNumber uint64      = 2
	defaultDistance      pb.Distance = pb.Distance_Dot
)

type QdrantConn struct {
	conn       *grpc.ClientConn
	clientColl pb.CollectionsClient
}

func (s *QdrantConn) Close() error {
	return s.conn.Close()
}

func (s *QdrantConn) NewColl(name string, vectorSize uint64) (vector.Coll, error) {
	req := pb.CreateCollection{
		CollectionName: name,
		VectorsConfig: &pb.VectorsConfig{Config: &pb.VectorsConfig_Params{
			Params: &pb.VectorParams{
				Size:     vectorSize,
				Distance: defaultDistance,
			},
		}},
		OptimizersConfig: &pb.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	}
	_, err := s.clientColl.Create(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return s.GetColl(name)
}

func (s *QdrantConn) GetColl(name string) (vector.Coll, error) {
	// check coll exist
	req := pb.GetCollectionInfoRequest{
		CollectionName: name,
	}
	_, err := s.clientColl.Get(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return &QdrantColl{client: pb.NewPointsClient(s.conn), nameColl: name}, nil
}

func (s *QdrantConn) DeleteColl(name string) error {
	req := pb.DeleteCollection{
		CollectionName: name,
	}
	_, err := s.clientColl.Delete(context.Background(), &req)
	if err != nil {
		return err
	}
	return nil
}

func NewQdrantConn(uri string) (*QdrantConn, error) {
	conn, err := grpc.DialContext(context.Background(), uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	qdrantClient := QdrantConn{
		conn:       conn,
		clientColl: pb.NewCollectionsClient(conn),
	}
	return &qdrantClient, nil
}
