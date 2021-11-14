package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(from, to string, offset, limit int64) error {
	src, err := os.OpenFile(from, os.O_RDONLY, 0)
	if err != nil {
		if os.IsExist(err) {
			panic("src file is not exist")
		}
		if os.IsPermission(err) {
			panic("src file not permitted")
		}
		panic(err)
	}

	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		if os.IsExist(err) {
			panic("file already exist")
		}
		panic(err)
	}

	defer dst.Close()

	info, err := src.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	if info.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("unable to set file offset: %v", err)
	}
	if limit == 0 {
		limit = info.Size() - offset
	}

	lReader := io.LimitReader(src, limit)
	pb := NewProgressBar(limit)
	pbWriter := pb.NewProgressBarWriter(dst)

	pb.Start()

	var written int64
	buf := make([]byte, 1)

	for {
		n, err := lReader.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("unable to read: %v", err)
		}
		if n == 0 {
			break
		}
		if _, err = pbWriter.Write(buf[:n]); err != nil {
			return fmt.Errorf("unable to write: %v", err)
		}
		written = written + int64(n)
	}

	pb.Finish()

	return nil
}
