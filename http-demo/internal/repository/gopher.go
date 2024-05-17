package repository

import (
	"context"
	"database/sql"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
)

// GopherRepository is the repository to handle the [model.Gopher] model database interactions.
type GopherRepository struct {
	mutex sync.Mutex
	db    *sql.DB
}

// NewGopherRepository returns a new [GopherRepository].
func NewGopherRepository(db *sql.DB) *GopherRepository {
	return &GopherRepository{
		db: db,
	}
}

// Find finds a gopher by id.
func (r *GopherRepository) Find(ctx context.Context, id int) (model.Gopher, error) {
	var gopher model.Gopher

	query, args, err := sq.
		Select("id", "name", "job").
		From("gophers").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return gopher, err
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	err = row.Scan(&gopher.ID, &gopher.Name, &gopher.Job)

	return gopher, err
}

// GopherRepositoryFindAllParams is a parameter for FindAll.
type GopherRepositoryFindAllParams struct {
	Name sql.NullString
	Job  sql.NullString
}

// FindAll finds all gophers.
func (r *GopherRepository) FindAll(ctx context.Context, params GopherRepositoryFindAllParams) ([]model.Gopher, error) {
	queryBuilder := sq.
		Select("id", "name", "job").
		From("gophers")

	if params.Name.Valid {
		queryBuilder.Where(sq.Eq{"name": params.Name.String})
	}

	if params.Job.Valid {
		queryBuilder.Where(sq.Eq{"name": params.Name.String})
	}

	queryBuilder.OrderBy("id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var gophers []model.Gopher
	for rows.Next() {
		var gopher model.Gopher
		if err = rows.Scan(&gopher.ID, &gopher.Name, &gopher.Job); err != nil {
			return nil, err
		}
		gophers = append(gophers, gopher)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gophers, nil
}

// GopherRepositoryCreateParams is a parameter for Create.
type GopherRepositoryCreateParams struct {
	Name string
	Job  sql.NullString
}

// Create creates a new gopher and returns its id.
func (r *GopherRepository) Create(ctx context.Context, params GopherRepositoryCreateParams) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query, args, err := sq.Insert("gophers").Columns("name", "job").Values(params.Name, params.Job).ToSql()

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Delete deletes an existing gopher by id.
func (r *GopherRepository) Delete(ctx context.Context, id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query, args, err := sq.Delete("gophers").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err
}
