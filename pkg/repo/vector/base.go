package vector

import (
	"strings"

	"github.com/wyubin/llm-svc/utils/repo/vector"
)

func NewConn(uri string) (vector.Conn, error) {
	protocol := strings.SplitN(uri, ":", 2)
	switch protocol[0] {
	case "qdrant":
		return NewQdrantConn(protocol[1])
	}
	return nil, nil
}
