package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"remindbot/lib/e"
	"remindbot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const (
	permissionCode = 0774
)

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(p *storage.Page) error {
	fp := filepath.Join(s.basePath, p.UserName)

	if err := os.MkdirAll(fp, permissionCode); err != nil {
		return e.Wrap("failed to create a dir", err)
	}

	fn, err := fileName(p)
	if err != nil {
		return e.Wrap("failed to hash", err)
	}

	fp = filepath.Join(fp, fn)

	file, err := os.Create(fp)
	if err != nil {
		return e.Wrap("failed to create a file", err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return e.Wrap("failed to encode gob", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	//joining needed path and opening the dir
	fp := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fp)
	if err != nil {
		return nil, e.Wrap("failed to access dir", err)
	}

	if len(files) == 0 {
		return nil, errors.New("You did not save anything yet")
	}

	//picking a random file and opening it
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]
	fp = filepath.Join(fp, file.Name())

	fileOpen, err := os.Open(fp)
	if err != nil {
		return nil, e.Wrap("failed to open file", err)
	}
	defer func() { _ = fileOpen.Close() }()

	//creating a page instance and decoding into it
	p := &storage.Page{}
	if err := gob.NewDecoder(fileOpen).Decode(&p); err != nil {
		return nil, e.Wrap("failed to decode gob", err)
	}

	return p, nil
}

func (s Storage) PickAll(userName string) ([]storage.Page, error) {
	//joining needed path and opening the dir
	fp := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fp)
	if err != nil {
		return nil, e.Wrap("failed to access dir", err)
	}

	if len(files) == 0 {
		return nil, errors.New("You did not save anything yet")
	}

	//ranging files and decoding
	p := &storage.Page{}
	var pSlice []storage.Page

	for _, f := range files {
		fp := filepath.Join(s.basePath, userName, f.Name())
		fileOpen, err := os.Open(fp)
		if err != nil {
			return nil, e.Wrap("failed to open a file", err)
		}

		if err := gob.NewDecoder(fileOpen).Decode(&p); err != nil {
			return nil, e.Wrap("failed to decode a file", err)
		}

		pSlice = append(pSlice, *p)
		fileOpen.Close()
	}

	return pSlice, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fm, err := fileName(p)
	if err != nil {
		return e.Wrap("failed to hash", err)
	}

	fp := filepath.Join(s.basePath, p.UserName, fm)

	if err := os.Remove(fp); err != nil {
		return e.Wrap("failed to delete a file or does not exist", err)
	}

	return nil
}

func (s Storage) IfExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
