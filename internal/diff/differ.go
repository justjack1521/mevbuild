package diff

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	errFailedCreateBinaryDiff = func(file File, err error) error {
		return fmt.Errorf("failed to create binary diff for file %s: %w", file.PatchFilePath, err)
	}
	errFailedCreateOrFindPath = func(err error) error {
		return fmt.Errorf("failed to create of find output path %w", err)
	}
)

type Differ struct {
	VerboseLevel int
	Timeout      time.Duration
}

func NewDiffer(verbose int, timeout int) *Differ {
	return &Differ{VerboseLevel: verbose, Timeout: time.Duration(timeout) * time.Second}
}

func (d *Differ) CreateBinaryDiff(file File) error {

	var dir = filepath.Dir(file.PatchFilePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errFailedCreateBinaryDiff(file, errFailedCreateOrFindPath(err))
	}

	if d.VerboseLevel > 0 {
		fmt.Println(fmt.Sprintf("directory created or already exists: %s", dir))

	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()

	var args = []string{
		"-j",
	}

	if d.VerboseLevel == 1 {
		args = append(args, "-v")
	}

	if d.VerboseLevel == 2 {
		args = append(args, "-vv")
	}

	args = append(args, file.OriginFilePath, file.NewFilePath, file.PatchFilePath)

	var cmd = exec.CommandContext(ctx, "C:\\jojodiff\\jojodiff.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if _, statErr := os.Stat(file.PatchFilePath); statErr == nil {
			if d.VerboseLevel > 0 {
				fmt.Println("patch file created despite non-zero exit code")
			}
		} else if ctx.Err() == context.DeadlineExceeded {
			return errFailedCreateBinaryDiff(file, fmt.Errorf("jojodiff command timed out"))
		} else {
			return errFailedCreateBinaryDiff(file, err)
		}
	} else {
		if d.VerboseLevel > 0 {
			fmt.Println(fmt.Sprintf("patch file created for %s", file.PatchFilePath))
		}
	}

	return nil

}

type File struct {
	OriginFilePath string
	NewFilePath    string
	PatchFilePath  string
}
