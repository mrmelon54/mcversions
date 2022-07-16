package utils

import (
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

var (
	ErrFileExists       = errors.New("file already exists")
	ErrCreateOutputFile = errors.New("failed to create output file")
	ErrStartDownload    = errors.New("failed to start download")
	ErrDuringDownload   = errors.New("failed during download")
	ErrDownloadSize     = errors.New("incorrect download size")
	ErrUnsafeDownload   = errors.New("sha1 hashes don't match... deleting the file for your safety")
	ErrUnsafeDelete     = errors.New("failed to delete unsafe file")
)

func DownloadJar(id *structure.PistonMetaId, dd structure.PistonMetaPackageDownloadsData) (int64, error) {
	filename := id.String() + "-" + path.Base(dd.URL)
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		return 0, ErrFileExists
	}
	out, err := os.Create(filename)
	if err != nil {
		return 0, ErrCreateOutputFile
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)
	resp, err := http.Get(dd.URL)
	if err != nil {
		return 0, ErrStartDownload
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	h := sha1.New()

	// Connect 'out' and 'h' as a single writer
	w := io.MultiWriter(out, h)

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		return 0, ErrDuringDownload
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	if n != dd.Size {
		return 0, ErrDownloadSize
	}

	sha1str := h.Sum(nil)
	if hex.EncodeToString(sha1str) != dd.Sha1 {
		err = os.Remove(filename)
		if err != nil {
			return 0, ErrUnsafeDelete
		}
		return 0, ErrUnsafeDownload
	}

	// File is safe
	return n, nil
}
