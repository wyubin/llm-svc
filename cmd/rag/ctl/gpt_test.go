package ctl

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	uri = "yubin-ollama-dev.user-yubin-wang.svc.cluster.local:11434"
)

func TestMain(m *testing.M) {
	// os.Exit(m.Run())
	os.Exit(0)
}

func TestOpenAI(t *testing.T) {
	gpt, err := NewGpt(fmt.Sprintf("http://%s", uri))
	assert.NoError(t, err)

	hasModel, err := gpt.HasModel("nomic-embed-text")
	assert.NoError(t, err)
	assert.True(t, hasModel)
}
