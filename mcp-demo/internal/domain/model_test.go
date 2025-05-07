package domain_test

import (
	"testing"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGopher(t *testing.T) {
	t.Parallel()

	gopher := &domain.Gopher{
		ID:   1,
		Name: "test name",
		Job:  "test job",
	}

	assert.Equal(t, int32(1), gopher.ID)
	assert.Equal(t, "test name", gopher.Name)
	assert.Equal(t, "test job", gopher.Job)
}
