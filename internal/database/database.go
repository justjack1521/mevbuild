package database

import (
	"github.com/justjack1521/mevpatch/internal/patch"

	"os"
	"path/filepath"
)

func CreateDatabaseFile(version patch.Version, configuration patch.Configuration) (string, error) {

	var path = filepath.Join(configuration.SourceOutputPath(), "patch", version.ToString(), "patching.sqlite")

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return "", err
	}

	db, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer db.Close()

	return path, nil

}
