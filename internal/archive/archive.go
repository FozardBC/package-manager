package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pm/internal/config"
	"regexp"
	"strings"
)

func CreateArchive(archivePath string, targets []config.Target) error {
	file, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("failed to create archive:%w", err)
	}
	defer file.Close()

	gWriter := gzip.NewWriter(file)
	defer gWriter.Close()
	tWriter := tar.NewWriter(gWriter)
	defer tWriter.Close()

	for _, target := range targets {
		matches, err := filepath.Glob(target.Path)
		if err != nil {
			return fmt.Errorf("failed to find mathes:%w", err)
		}

		exclude, err := compileExcludes(target.Exclude)
		if err != nil {
			return fmt.Errorf("failed to compile excludes: %w", err)
		}

		for _, match := range matches {
			filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				rel, err := filepath.Rel(".", path)
				if err != nil {
					return err
				}

				if matchesAny(rel, exclude) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				header, err := tar.FileInfoHeader(info, "")
				if err != nil {
					return err
				}

				header.Name = filepath.ToSlash(rel)

				if err := tWriter.WriteHeader(header); err != nil {
					return err

				}

				if info.Mode().IsRegular() {
					file, err := os.Open(path)
					if err != nil {
						return err
					}
					defer file.Close()
					io.Copy(tWriter, file)

				}
				return nil
			},
			)
		}
	}
	return nil
}

func ExtractArchive(archivePath, dest string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}

			out, err := os.Create(target)
			if err != nil {
				return err
			}

			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}

			if err := out.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func compileExcludes(pattern string) (*regexp.Regexp, error) {
	var result *regexp.Regexp

	shielded := regexp.QuoteMeta(pattern)
	regex := strings.ReplaceAll(shielded, `\*`, `.*`)

	result, err := regexp.Compile("^" + regex + "$")
	if err != nil {

		return nil, fmt.Errorf("failed to compile exclude pattern '%s': %w", pattern, err)
	}

	return result, nil
}

func matchesAny(s string, r *regexp.Regexp) bool {

	if r.MatchString(s) {
		return true
	}

	return false

}
