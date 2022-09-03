package file

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"valuation/pkg/convertx"
)

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	p, _ := filepath.Abs(os.Args[0])
	return p
}

// RealPath get absolute filepath, based on built executable file
func RealPath(fp string) (string, error) {
	if path.IsAbs(fp) {
		return fp, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, fp), err
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// Basename get filepath base name
func Basename(fp string) string {
	return path.Base(fp)
}

// Dir get filepath dir name
func Dir(fp string) string {
	return path.Dir(fp)
}

func InsureDir(fp string) error {
	if IsExist(fp) {
		return nil
	}
	return os.MkdirAll(fp, os.ModePerm)
}

// EnsureDir mkdir dir if not exist
func EnsureDir(fp string) error {
	return os.MkdirAll(fp, os.ModePerm)
}

// EnsureDirRW ensure the dataDir and make sure it's rw-able
func EnsureDirRW(dataDir string) error {
	if err := EnsureDir(dataDir); err != nil {
		return err
	}

	checkFile := fmt.Sprintf("%s/rw.%d", dataDir, time.Now().UnixNano())
	fd, err := Create(checkFile)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("open %s: rw permission denied", dataDir)
		}
		return err
	}
	Close(fd)

	if err = Remove(checkFile); err != nil {
		return err
	}

	return nil
}

// Create  one file
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Remove one file
func Remove(name string) error {
	return os.Remove(name)
}

func Close(fd *os.File) error {
	return fd.Close()
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(fp string) bool {
	f, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

// Unlink delete file
func Unlink(fp string) error {
	if !IsExist(fp) {
		return nil
	}
	return os.Remove(fp)
}

// FileMTime get file modified time
func FileMTime(fp string) (int64, error) {
	f, err := os.Stat(fp)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}

// FileSize get file size as how many bytes
func FileSize(fp string) (int64, error) {
	f, err := os.Stat(fp)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

// DirsUnder list dirs under dirPath
func DirsUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}

	return ret, nil
}

// FilesUnder list files under dirPath
func FilesUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if !fs[i].IsDir() {
			ret = append(ret, fs[i].Name())
		}
	}

	return ret, nil
}

func ReadBytes(cpath string) ([]byte, error) {
	if !IsExist(cpath) {
		return nil, fmt.Errorf("%s not exists", cpath)
	}

	if !IsFile(cpath) {
		return nil, fmt.Errorf("%s not file", cpath)
	}

	return os.ReadFile(cpath)
}

func ReadString(cpath string) (string, error) {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func ReadStringTrim(cpath string) (string, error) {
	out, err := ReadString(cpath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func WriteBytes(filePath string, b []byte) (int, error) {
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return 0, err
	}
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

func WriteString(filePath string, s string) (int, error) {
	return WriteBytes(filePath, convertx.String2Bytes(s))
}

func Download(toFile, url string) error {
	out, err := os.Create(toFile)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.Body == nil {
		return fmt.Errorf("%s response body is nil", url)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func MD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
