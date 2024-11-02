package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/justjack1521/mevpatch/internal/manifest"
	"github.com/justjack1521/mevpatch/internal/patch"
	"os"
	"path/filepath"
)

func main() {

	var p string
	var v string
	var n int

	flag.StringVar(&p, "p", "", "patch profile name")
	flag.StringVar(&v, "v", "", "current patch target version")
	flag.IntVar(&n, "n", 5, "number of historic versions to include")
	flag.Parse()

	configuration, err := patch.NewConfiguration(p)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Remote host set as %s", configuration.Host))
	fmt.Println(fmt.Sprintf("Source path set as %s", configuration.Source))

	if err := configuration.Test(); err != nil {
		panic(err)
	}

	fmt.Println("Configuration test passed successfully")

	version, err := patch.NewVersion(v)
	if err != nil {
		panic(err)
	}

	var previous = version.GeneratePreviousVersions(n)

	fmt.Println(fmt.Sprintf("Current version set as: %s", version.ToString()))
	fmt.Println(fmt.Sprintf("New Minor Version: %v", version.IsNewMinorVersion()))
	fmt.Println(fmt.Sprintf("Patching for %d previous versions", len(previous)))

	for _, prev := range previous {
		fmt.Println(fmt.Sprintf("- %s", prev.ToString()))
	}

	ctx, err := patch.NewContext(version, configuration)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%d files found in version %s", len(ctx.Files), ctx.Version.ToString()))

	if err := ctx.MountPrePatchFiles(previous); err != nil {
		fmt.Println(err)
	}

	if err := ctx.CreatePatchFiles(); err != nil {
		fmt.Println(err)
	}

	var m = &manifest.Manifest{Version: ctx.Version.ToString()}

	m.Files = make([]*manifest.File, len(ctx.Files))

	for i, file := range ctx.Files {
		m.Files[i] = manifest.NewFile(file)
	}

	file, err := os.Create(filepath.Join(ctx.Configuration.SourceOutputPath(), "manifest.json"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var encoder = json.NewEncoder(file)
	encoder.SetIndent("", "	")
	if err := encoder.Encode(m); err != nil {
		panic(err)
	}

}
