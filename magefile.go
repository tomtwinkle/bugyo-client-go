// +build mage

package bugyoclient

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuildEnv struct {
	OS           string
	ARCH         string
	File         string
	Compress     string
	CompressFile string
}

func Build() error {
	mg.Deps(InstallDeps)

	var buildEnvList = []BuildEnv{
		{"linux", "386", "bugyoclient", "gz", "bugyoclient-%s.rev-%s-i386_linux.gz"},
		{"linux", "amd64", "bugyoclient", "gz", "bugyoclient-%s.rev-%s-arm64_linux.gz"},
		{"darwin", "386", "bugyoclient", "gz", "bugyoclient-%s.rev-%s-i386_mac.gz"},
		{"darwin", "amd64", "bugyoclient", "gz", "bugyoclient-%s.rev-%s-arm64_mac.gz"},
		{"windows", "386", "bugyoclient.exe", "zip", "bugyoclient-%s.rev-%s_win32.zip"},
		{"windows", "amd64", "bugyoclient.exe", "zip", "bugyoclient-%s.rev-%s_win64.zip"},
	}

	v, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		fmt.Printf("%v\n", err)
		return errors.WithStack(err)
	}
	version := string(v)
	version = strings.Trim(version, "\r\n")

	r, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		fmt.Printf("%v\n", err)
		return errors.WithStack(err)
	}
	revision := string(r)
	revision = strings.Trim(revision, "\r\n")

	fmt.Println(fmt.Sprintf("Building %s.rev-%s ...", version, revision))
	for _, e := range buildEnvList {
		filePath := filepath.Join("builds", e.File)
		ldflags := fmt.Sprintf("-s -w -X main.version=%s -X main.revision=%s", version, revision)

		cmd := exec.Command("go", "build", "-o", filePath, "-ldflags", ldflags, ".")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("GOOS=%s", e.OS),
			fmt.Sprintf("GOARCH=%s", e.ARCH),
		)
		if err := cmd.Run(); err != nil {
			fmt.Printf("%v\n", err)
			return errors.WithStack(err)
		}
		compressFilePath := filepath.Join("builds", fmt.Sprintf(e.CompressFile, version, revision))
		if e.Compress == "gz" {
			if err := makeGzip(compressFilePath, filePath); err != nil {
				fmt.Printf("%v\n", err)
				return errors.WithStack(err)
			}
		} else if e.Compress == "zip" {
			if err := makeZip(compressFilePath, filePath); err != nil {
				fmt.Printf("%v\n", err)
				return errors.WithStack(err)
			}
		}
		os.Remove(filePath)
	}
	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	return nil
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("builds")
}

func makeGzip(compressFilePath string, filePath string) error {
	fmt.Printf("make compress %s\n", compressFilePath)
	orgfile, err := ioutil.ReadFile(filePath)

	gfile, err := os.Create(compressFilePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer gfile.Close()
	zw, err := gzip.NewWriterLevel(gfile, gzip.BestCompression)
	if err != nil {
		return errors.WithStack(err)
	}
	defer zw.Close()

	if _, err := zw.Write(orgfile); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func makeZip(compressFilePath string, filePath string) error {
	fmt.Printf("make compress %s\n", compressFilePath)
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	newZipFile, err := os.Create(compressFilePath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return errors.WithStack(err)
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
