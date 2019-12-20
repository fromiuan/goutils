package tools

import (
	"bufio"
	"crypto/md5"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Md5File get file md5
func Md5File(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	h := md5.New()
	io.Copy(h, file)
	return string(h.Sum(nil))
}

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// SearchFile Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	isLongLine := false
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// RemoveFile remove file
func RemoveFile(file string) error {
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}

// ReadFileSQL read sql file
// readfile("test.sql")
func ReadFileSQL(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	defer file.Close()

	var sqls []string
	var tmpline string
	scanner := bufio.NewScanner(file)
	commentBegin := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, "*/") {
			commentBegin = false
			continue
		}
		if len(line) > 0 {
			if strings.HasPrefix(line, "//") {
				continue
			} else if strings.HasPrefix(line, "--") {
				continue
			} else if strings.HasPrefix(line, "/*") {
				if strings.HasSuffix(line, "*/") {
					continue
				}
				commentBegin = true
				continue
			} else {
				if commentBegin {
					continue
				} else {
					// TODO
					if strings.HasSuffix(line, ";") {
						sqls = append(sqls, tmpline+line)
						tmpline = ""
					} else {
						tmpline = tmpline + line
					}
				}
			}
		}
	}
	return sqls
}
