package manifest

import (
	"encoding/json"
	"fmt"
	mevmanifest "github.com/justjack1521/mevmanifest/pkg/genproto"
	"github.com/justjack1521/mevpatch/internal/patch"
	"os"
	"path/filepath"
)

func CreateManifestFile(ctx *patch.Context, bundles []patch.BundleFile) error {
	var m = &mevmanifest.Manifest{Version: ctx.Version.ToString()}

	m.Files = make([]*mevmanifest.File, len(ctx.Files))
	m.Bundles = make([]*mevmanifest.Bundle, len(bundles))

	for i, file := range ctx.Files {
		m.Files[i] = NewFile(file)
	}

	for i, bundle := range bundles {
		m.Bundles[i] = NewBundle(bundle)
	}

	var path = filepath.Join(ctx.Configuration.SourceOutputPath(), "patch", ctx.Version.ToString(), fmt.Sprintf("%s_manifest.json", ctx.Version.ToString()))

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var encoder = json.NewEncoder(file)
	encoder.SetIndent("", "	")
	if err := encoder.Encode(m); err != nil {
		return err
	}
	return err
}
