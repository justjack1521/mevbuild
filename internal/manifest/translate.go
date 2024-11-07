package manifest

import (
	mevmanifest "github.com/justjack1521/mevmanifest/pkg/genproto"
	"github.com/justjack1521/mevpatch/internal/patch"
)

func NewFile(file *patch.File) *mevmanifest.File {

	var result = &mevmanifest.File{
		Id:           file.ID.String(),
		Path:         file.NormalPath,
		Checksum:     file.Checksum,
		Size:         file.Size,
		Patches:      make([]*mevmanifest.PatchFile, len(file.OutputFiles)),
		DownloadPath: file.DownloadPath,
	}

	for index, out := range file.OutputFiles {
		result.Patches[index] = NewPatchFile(out)
	}

	return result

}

func NewPatchFile(file patch.OutputPatchFile) *mevmanifest.PatchFile {
	return &mevmanifest.PatchFile{
		Version:  file.Version.ToString(),
		Checksum: file.Checksum,
		Size:     file.Size,
	}
}

func NewBundle(file patch.BundleFile) *mevmanifest.Bundle {
	return &mevmanifest.Bundle{
		Version:      file.Version.ToString(),
		DownloadPath: file.DownloadPath,
		Size:         file.Size,
		Checksum:     file.Checksum,
	}
}
