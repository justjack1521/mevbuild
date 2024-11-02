package patch

import (
	"errors"
	"fmt"
	"github.com/justjack1521/mevpatch/internal/diff"
	"io/fs"
	"os"
	"path/filepath"
)

type Context struct {
	Version       Version
	Configuration Configuration
	Files         []*File
}

func NewContext(version Version, configuration Configuration) (*Context, error) {

	var paths = make([]string, 0)

	var current = filepath.Join(configuration.SourceInputPath(), version.ToString())

	if err := filepath.WalkDir(current, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() == false {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	var ctx = &Context{
		Configuration: configuration,
		Version:       version,
		Files:         make([]*File, len(paths)),
	}

	for index, path := range paths {
		file, err := NewFile(ctx, path)
		if err != nil {
			return nil, err
		}
		ctx.Files[index] = file
	}

	return ctx, nil

}

func (f *File) CreatePathFilePath(c Configuration, input InputPatchFile) string {
	return filepath.Join(c.SourceOutputPath(), "patch", f.NormalPath, fmt.Sprintf("%s_%s.%s", input.Version.ToString(), f.Version.ToString(), c.Suffix))
}

func (c *Context) CreatePatchFiles() error {

	var differ = diff.NewDiffer(c.Configuration.Differ.VerboseLevel, c.Configuration.Differ.Timeout)

	for _, file := range c.Files {
		for _, patch := range file.InputFiles {

			var input = diff.File{
				OriginFilePath: patch.LocalPath,
				NewFilePath:    file.LocalPath,
				PatchFilePath:  patch.PatchFilePath,
			}

			if err := differ.CreateBinaryDiff(input); err != nil {
				return err
			}

			output, err := NewOutputPatchFile(c, patch)
			if err != nil {
				return err
			}

			file.OutputFiles = append(file.OutputFiles, output)

		}
	}

	return nil

}

func (c *Context) MountPrePatchFiles(versions []Version) error {

	var counter int
	var skipped int

	for _, file := range c.Files {
		for _, version := range versions {
			previous, err := file.FindPreviousFileVersion(c.Configuration, version)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					skipped++
					continue
				}
				return err
			}
			file.InputFiles = append(file.InputFiles, previous)
			counter++
		}
	}

	fmt.Println(fmt.Sprintf("%d files skipped for patching across %d versions", skipped, len(versions)))
	fmt.Println(fmt.Sprintf("%d files mounted for patching across %d versions", counter, len(versions)))

	return nil

}
