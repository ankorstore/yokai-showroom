package domain_test

import (
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGopher(t *testing.T) {
	t.Parallel()

	gopher := &domain.Gopher{
		ID:   1,
		Name: "test",
		Job:  sql.NullString{String: "test", Valid: true},
	}

	assert.Equal(t, int32(1), gopher.ID)
	assert.Equal(t, "test", gopher.Name)
	assert.Equal(t, "test", gopher.Job.String)
}
