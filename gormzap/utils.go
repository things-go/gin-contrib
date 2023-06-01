package gormzap

import (
	"runtime"
	"strconv"
	"strings"
)

var (
	gormPackage    = "gorm.io/gorm"
	gormzapPackage = "github.com/things-go/gin-contrib/gormzap"
)

// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum(skipPackages ...string) string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !skipPackage(file, skipPackages...) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

func skipPackage(file string, skipPackages ...string) bool {
	if strings.HasSuffix(file, "_test.go") {
		return false
	}
	if strings.Contains(file, gormPackage) ||
		strings.Contains(file, gormzapPackage) {
		return true
	}
	for _, p := range skipPackages {
		if strings.Contains(file, p) {
			return true
		}
	}
	return false
}
