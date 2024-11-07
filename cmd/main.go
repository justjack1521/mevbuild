package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/justjack1521/mevpatch/internal/database"
	"github.com/justjack1521/mevpatch/internal/manifest"
	"github.com/justjack1521/mevpatch/internal/patch"
	_ "modernc.org/sqlite"
)

func main() {

	var t string
	var v string
	var n int

	flag.StringVar(&t, "t", "", "target application name")
	flag.StringVar(&v, "v", "", "current patch target version")
	flag.IntVar(&n, "n", 5, "number of historic versions to include")
	flag.Parse()

	configuration, err := patch.NewConfiguration(t)
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

	ctx, err := patch.NewContext(configuration, version, previous)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%d files found in version %s", len(ctx.Files), ctx.Version.ToString()))

	if err := ctx.MountPrePatchFiles(); err != nil {
		panic(err)
	}

	if err := ctx.CreatePatchFiles(); err != nil {
		panic(err)
	}

	var bundler = ctx.NewBundler()

	bundles, err := bundler.BuildPatchFiles()
	if err != nil {
		panic(err)
	}

	if err := manifest.CreateManifestFile(ctx, bundles); err != nil {
		panic(err)
	}

	path, err := database.CreateDatabaseFile(ctx.Version, ctx.Configuration)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var repository = database.NewPatchFileRepository(db)
	if err := repository.Initialise(); err != nil {
		panic(err)
	}

	if err := repository.CreateApplicationVersion(t, version); err != nil {
		panic(err)
	}

	for _, file := range ctx.Files {
		if err := repository.CreateApplicationFile(t, file); err != nil {
			panic(err)
		}
	}

}
