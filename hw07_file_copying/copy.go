package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	WithProgressBar = true

	ErrSrcFileIsNotExist     = errors.New("src file is not exist")
	ErrSrcFileIsNotPermitted = errors.New("src file is not permitted")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(srcFilePath, dstFilePath string, offset, limit int64) error {
	src, err := os.OpenFile(srcFilePath, os.O_RDONLY, 0)
	if err != nil {
		if !os.IsExist(err) {
			return ErrSrcFileIsNotExist
		}
		if !os.IsPermission(err) {
			return ErrSrcFileIsNotPermitted
		}
		return fmt.Errorf("unknown error: %v", err)
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file info: %v", err)
	}
	size := info.Size()
	if size == 0 {
		return ErrUnsupportedFile
	}
	if size <= offset {
		return ErrOffsetExceedsFileSize
	}

	dst, err := os.Create(dstFilePath)
	if err != nil {
		return fmt.Errorf("unable to create dst file: %v", err)
	}
	defer dst.Close()

	if _, err = src.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("unable to set file offset: %v", err)
	}

	total := size - offset
	if limit == 0 || limit > total {
		limit = total
	} else {
		total = limit
	}

	if WithProgressBar {
		return pbCopyN(dst, io.LimitReader(src, limit), total)
	} else {
		return copyN(dst, io.LimitReader(src, limit), total)
	}
}

func pbCopyN(dst io.Writer, src io.Reader, n int64) error {
	pb := NewProgressBar(n)
	pb.Start()
	err := copyN(pb.NewProgressBarWriter(dst), src, n)
	pb.Finish()

	return err
}

func copyN(dst io.Writer, src io.Reader, n int64) error {
	if _, err := io.CopyN(dst, src, n); err != nil && err != io.EOF {
		return fmt.Errorf("not copied: %v", err)
	}
	return nil
}
