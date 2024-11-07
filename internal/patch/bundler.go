package patch

import (
	"archive/zip"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type Bundler struct {
	Configuration Configuration
	Target        Version
	Patches       map[Version][]OutputPatchFile
}

type BundleFile struct {
	Version      Version
	DownloadPath string
	Size         int64
	Checksum     string
}

func (b *Bundler) BuildPatchFiles() ([]BundleFile, error) {

	var files = make([]BundleFile, 0)

	for key, value := range b.Patches {
		file, err := b.BundlePatchVersionFiles(key, b.Target, value)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (b *Bundler) BundlePatchVersionFiles(from Version, to Version, files []OutputPatchFile) (BundleFile, error) {

	var file = fmt.Sprintf("%s_patch.bin", from.ToString())

	download, err := url.JoinPath(b.Configuration.Host, "downloads", b.Configuration.AppName, "patch", to.ToString(), file)
	if err != nil {
		return BundleFile{}, err
	}

	var path = filepath.Join(b.Configuration.SourceOutputPath(), "patch", to.ToString(), file)

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return BundleFile{}, fmt.Errorf("failed to create directory for patch file bundle: %w", err)
	}

	bundle, err := os.Create(path)
	if err != nil {
		return BundleFile{}, fmt.Errorf("failed to create patch file bundle: %w", err)
	}
	defer bundle.Close()

	var writer = zip.NewWriter(bundle)
	defer writer.Close()

	for _, f := range files {
		_, err := b.BundlePatchFile(b.Configuration, writer, f)
		if err != nil {
			return BundleFile{}, err
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return BundleFile{}, err
	}

	checksum, err := GetChecksum(path)
	if err != nil {
		return BundleFile{}, err
	}

	fmt.Println(fmt.Sprintf("Version bundle created at: %s", path))

	return BundleFile{
		Version:      from,
		DownloadPath: download,
		Size:         info.Size(),
		Checksum:     checksum,
	}, nil

}

func (b *Bundler) BundlePatchFile(c Configuration, w *zip.Writer, output OutputPatchFile) (int64, error) {

	fmt.Println(fmt.Sprintf("Bundling: %s", output.PatchFileTempPath))

	file, err := os.Open(output.PatchFileTempPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return 0, err
	}
	header.Name = fmt.Sprintf("%s.%s", output.ID.String(), c.Suffix)
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return 0, err
	}

	count, err := io.Copy(writer, file)
	if err != nil {
		return 0, err
	}

	return count, nil

}
