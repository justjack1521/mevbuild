package manifest

import (
	"github.com/justjack1521/mevpatch/internal/patch"
)

func NewFile(file *patch.File) *File {

	var result = &File{
		Path:         file.NormalPath,
		Checksum:     file.Checksum,
		Patches:      make([]*PatchFile, len(file.OutputFiles)),
		DownloadPath: file.DownloadPath,
	}

	for index, out := range file.OutputFiles {
		result.Patches[index] = NewPatchFile(out)
	}

	return result

}

func NewPatchFile(file patch.OutputPatchFile) *PatchFile {
	return &PatchFile{
		Version:  file.Version.ToString(),
		Checksum: file.Checksum,
		Size:     int32(file.Size),
	}
}
