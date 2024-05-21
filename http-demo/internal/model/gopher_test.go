package model_test

import (
	"database/sql"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGopher(t *testing.T) {
	t.Parallel()

	gopher := &model.Gopher{
		ID:   1,
		Name: "test",
		Job:  sql.NullString{String: "test", Valid: true},
	}

	assert.Equal(t, int32(1), gopher.ID)
	assert.Equal(t, "test", gopher.Name)
	assert.Equal(t, "test", gopher.Job.String)
}
