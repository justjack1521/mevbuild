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
	Previous      []Version
	Configuration Configuration
	Files         []*File
}

func NewContext(configuration Configuration, target Version, previous []Version) (*Context, error) {

	var paths = make([]string, 0)

	var current = filepath.Join(configuration.SourceInputPath(), target.ToString())

	if err := filepath.WalkDir(current, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() == false {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	var actual = make([]Version, 0)
	for _, pre := range previous {
		var path = filepath.Join(configuration.SourceInputPath(), pre.ToString())
		if _, err := os.Stat(path); err != nil {
			continue
		}
		actual = append(actual, pre)
	}

	var ctx = &Context{
		Configuration: configuration,
		Version:       target,
		Previous:      actual,
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

func (f *File) CreatePatchFileName(c Configuration, input InputPatchFile) string {
	return fmt.Sprintf("%s_%s_%s.%s", f.NormalPath, input.Version.ToString(), f.Version.ToString(), c.Suffix)
}

func (f *File) CreatePatchFilePath(c Configuration, input InputPatchFile) string {
	return filepath.Join(c.SourceOutputPath(), "patch", "temp", f.Version.ToString(), f.NormalPath, f.CreatePatchFileName(c, input))
}

func (c *Context) NewBundler() *Bundler {

	var bundler = &Bundler{
		Configuration: c.Configuration,
		Target:        c.Version,
		Patches:       make(map[Version][]OutputPatchFile),
	}

	for _, previous := range c.Previous {
		bundler.Patches[previous] = make([]OutputPatchFile, 0)
	}

	for _, file := range c.Files {
		for _, p := range file.OutputFiles {
			bundler.Patches[p.Version] = append(bundler.Patches[p.Version], p)
		}
	}

	return bundler

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

func (c *Context) MountPrePatchFiles() error {

	var counter int
	var skipped int

	for _, file := range c.Files {
		for _, version := range c.Previous {
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

	fmt.Println(fmt.Sprintf("%d files skipped for patching across %d versions", skipped, len(c.Previous)))
	fmt.Println(fmt.Sprintf("%d files mounted for patching across %d versions", counter, len(c.Previous)))

	return nil

}
