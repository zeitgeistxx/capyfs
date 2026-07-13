package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "minecraft"

var DeafultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := range sliceLen {
		from, to := (i * blockSize), (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) AbsolutePath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

type StoreOpts struct {
	Root string
	PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DeafultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	absolutePath := fmt.Sprintf("%s/%s", s.Root, pathKey.AbsolutePath())

	_, err := os.Stat(absolutePath)
	if err == fs.ErrNotExist {
		return false
	}
	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()

	firstAbsolutePath := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())

	return os.RemoveAll(firstAbsolutePath)
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	absolutePath := fmt.Sprintf("%s/%s", s.Root, pathKey.AbsolutePath())
	return os.Open(absolutePath)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	absoluteFilePath := fmt.Sprintf("%s/%s", s.Root, pathKey.AbsolutePath())

	f, err := os.Create(absoluteFilePath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, absoluteFilePath)

	return nil
}
