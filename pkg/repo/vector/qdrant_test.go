package vector

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wyubin/go-utils/repo/vector"
)

const (
	uri       = "yubin-qdrant-dev.user-yubin-wang.svc.cluster.local:6334"
	lenVector = 128
)

func TestMain(m *testing.M) {
	// os.Exit(m.Run())
	os.Exit(0)
}

func TestQdrantConn(t *testing.T) {
	conn, err := NewQdrantConn(uri)
	assert.NoError(t, err)
	defer conn.Close()
	var coll vector.Coll
	_, err = conn.GetColl("test")
	if err == nil {
		err = conn.DeleteColl("test")
		assert.NoError(t, err)
	}
	coll, err = conn.NewColl("test", lenVector)
	assert.NoError(t, err)
	fmt.Printf("coll: %+v\n", coll)

	// final clear
	conn.DeleteColl("test")
}

func TestQdrantColl(t *testing.T) {
	conn, err := NewQdrantConn(uri)
	assert.NoError(t, err)
	defer conn.Close()
	var coll vector.Coll
	coll, _ = conn.NewColl("test", 2)
	uid1, err := coll.InsertOne([]float32{0.2, 0.2}, map[string]string{"type": "circle"})
	assert.NoError(t, err)
	fmt.Printf("uid1: %s\n", uid1)
	uid2, err := coll.InsertOne([]float32{0.9, 0.8}, map[string]string{"type": "rectangle"})
	assert.NoError(t, err)
	fmt.Printf("uid2: %s\n", uid2)
	points := []vector.Point{}

	err = coll.VectorSearch([]float32{0.1, 0.8}, &points, map[string]interface{}{"limit": 10})
	assert.NoError(t, err)
	for _, p := range points {
		fmt.Printf("p[meta]: %+v, p[score]: %+v\n", p.GetMeta(), p.GetScore())
	}
	fmt.Printf("points: %+v\n", points[0].GetMeta())
	// final clear
	conn.DeleteColl("test")
}
