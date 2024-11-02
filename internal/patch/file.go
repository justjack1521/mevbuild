package patch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	errVersionMarkerNotInPath = func(path string, marker string) error {
		return fmt.Errorf("version marker %s not found in path %s", path, marker)
	}
)

type File struct {
	Version      Version
	LocalPath    string
	NormalPath   string
	DownloadPath string
	Checksum     string
	Size         int64
	InputFiles   []InputPatchFile
	OutputFiles  []OutputPatchFile
}

type InputPatchFile struct {
	Version       Version
	LocalPath     string
	NormalPath    string
	PatchFilePath string
}

type OutputPatchFile struct {
	Version       Version
	PatchFilePath string
	Checksum      string
	DownloadPath  string
	Size          int64
}

func (x OutputPatchFile) ToString() string {
	return fmt.Sprintf("[Path: %s] [Version: %s] [Checksum: %s] [Size: %d]", x.PatchFilePath, x.Version.ToString(), x.Checksum, x.Size)
}

func NewOutputPatchFile(ctx *Context, input InputPatchFile) (OutputPatchFile, error) {

	file, err := os.Open(input.PatchFilePath)
	if err != nil {
		return OutputPatchFile{}, err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return OutputPatchFile{}, err
	}

	var hash = sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return OutputPatchFile{}, err
	}

	var checksum = hex.EncodeToString(hash.Sum(nil))

	download, err := url.JoinPath(ctx.Configuration.Host, "downloads", ctx.Configuration.AppName, "patch", input.NormalPath, file.Name())
	if err != nil {
		return OutputPatchFile{}, err
	}

	return OutputPatchFile{
		Version:       input.Version,
		PatchFilePath: input.PatchFilePath,
		DownloadPath:  download,
		Checksum:      checksum,
		Size:          stats.Size(),
	}, nil

}

func (f *File) ToString() string {
	return fmt.Sprintf("[Local Path: %s] [Normal Path: %s] [Size: %d] [Checksum: %s]", f.LocalPath, f.NormalPath, f.Size, f.Checksum)
}

func NewFile(ctx *Context, path string) (*File, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var hash = sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	var checksum = hex.EncodeToString(hash.Sum(nil))

	var normalised = filepath.ToSlash(path)
	var marker = ctx.Version.ToString()

	var index = strings.Index(normalised, marker)
	if index == -1 {
		return nil, errVersionMarkerNotInPath(path, marker)
	}
	var start = index + len(marker) + 1
	var normal = normalised[start:]

	download, err := url.JoinPath(ctx.Configuration.Host, "downloads", ctx.Configuration.AppName, "source", normal)
	if err != nil {
		return nil, err
	}

	return &File{
		Version:      ctx.Version,
		LocalPath:    path,
		NormalPath:   normal,
		DownloadPath: download,
		Checksum:     checksum,
		Size:         stats.Size(),
	}, nil

}

func (f *File) FindPreviousFileVersion(configuration Configuration, version Version) (InputPatchFile, error) {
	var path = filepath.Join(configuration.SourceInputPath(), version.ToString(), f.NormalPath)
	_, err := os.Stat(path)
	if err != nil {
		return InputPatchFile{}, err
	}
	var output = InputPatchFile{
		Version:   version,
		LocalPath: path,
	}

	output.PatchFilePath = f.CreatePathFilePath(configuration, output)

	return output, nil

}
