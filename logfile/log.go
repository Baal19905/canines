package logfile

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"time"
)

type Logfile struct {
	file *os.File
	date string
}

func SetLogfile(path string) *Logfile {
	dir, _ := filepath.Split(path)
	if exists, err := pathExists(dir); err != nil && exists == false {
		os.MkdirAll(dir, 0755)
	} else if err != nil {
		fmt.Printf("pathExists err: %s\n", err.Error())
		return nil
	}
	nowDate := time.Now().Format("20060102")
	truePath := path + "." + nowDate
	logfile, err := os.OpenFile(truePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("OpenFile err: %s\n", err.Error())
		return nil
	}
	err = os.Symlink(truePath, path)
	if err != nil {
		fmt.Printf("Symlink err: %s\n", err.Error())
	}
	return &Logfile{
		file: logfile,
		date: nowDate,
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (lf *Logfile) UpdateLogfile() {
	nowDate := time.Now().Format("20060102")
	if lf.date != nowDate {
		path := lf.file.Name()
		newPath := path + "." + nowDate
		newFile, err1 := os.OpenFile(newPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		err2 := os.Symlink(newPath, path)
		if err1 != nil && err2 != nil {
			lf.file.Close()
			lf.file = newFile
		}
	}
}

func (lf *Logfile) Write(p []byte) (n int, err error) {
	if lf.file == nil {
		return 0, errors.New("file not open")
	}
	lf.UpdateLogfile()
	return lf.file.Write(p)
}
