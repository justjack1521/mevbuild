package database

import (
	"context"
	"database/sql"
	mevmanifest "github.com/justjack1521/mevmanifest/pkg/gensql"
	"github.com/justjack1521/mevpatch/internal/patch"
)

type PatchFileRepository struct {
	queries *mevmanifest.Queries
}

func NewPatchFileRepository(db *sql.DB) *PatchFileRepository {
	return &PatchFileRepository{queries: mevmanifest.New(db)}
}

func (r *PatchFileRepository) Initialise() error {
	if err := r.queries.CreateApplicationVersionTable(context.Background()); err != nil {
		return err
	}
	if err := r.queries.CreateApplicationFileTable(context.Background()); err != nil {
		return err
	}
	return nil
}

func (r *PatchFileRepository) CreateApplicationVersion(app string, version patch.Version) error {
	var args = mevmanifest.CreateApplicationVersionParams{
		Name:  app,
		Major: int64(version.Major),
		Minor: int64(version.Minor),
		Patch: int64(version.Patch),
	}
	return r.queries.CreateApplicationVersion(context.Background(), args)
}

func (r *PatchFileRepository) CreateApplicationFile(app string, file *patch.File) error {
	var args = mevmanifest.CreateApplicationFileParams{
		Path:        file.NormalPath,
		Size:        file.Size,
		Timestamp:   file.LastModified.Unix(),
		Application: app,
	}
	return r.queries.CreateApplicationFile(context.Background(), args)
}
