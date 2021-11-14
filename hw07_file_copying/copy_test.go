package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	WithProgressBar = true
	t.Run("unsupported file to copy", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/urandom.txt", 0, 0)
		defer os.Remove("/tmp/urandom.txt")
		require.EqualError(t, err, ErrUnsupportedFile.Error(), "/dev/urandom is unsupported file")

		err = Copy("/dev/random", "/tmp/random.txt", 0, 0)
		defer os.Remove("/tmp/random.txt")
		require.EqualError(t, ErrUnsupportedFile, err.Error(), "/dev/random is unsupported file")
	})

	t.Run("not exist src file", func(t *testing.T) {
		err := Copy("/testdata/nonexist.txt", "/tmp/nonexist.txt", 0, 0)
		require.EqualError(t, err, ErrSrcFileIsNotExist.Error(), "must be error about file is not exist")
	})

	t.Run("exist dst file", func(t *testing.T) {
		srcFilePath := "testdata/input.txt"
		dstFilePath := "/tmp/output.txt"

		dstFile, err := os.Create(dstFilePath)
		require.NoError(t, err, "can't create dst file")
		defer os.Remove(dstFilePath)
		defer dstFile.Close()

		err = Copy(srcFilePath, dstFilePath, 0, 0)
		require.EqualError(t, err, ErrDstFileAlreadyExists.Error(), "must be error if dst file is already exist")
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		srcFilePath := "testdata/input.txt"
		dstFilePath := "/tmp/output.txt"

		srcFile, err := os.OpenFile(srcFilePath, os.O_RDONLY, 0)
		require.NoError(t, err, "unable to open src file")
		defer srcFile.Close()

		info, _ := srcFile.Stat()
		err = Copy(srcFilePath, dstFilePath, info.Size()+10, 0)
		if err == nil {
			defer os.Remove(dstFilePath)
		}
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error(), "must be error if offset exceeds file size")
	})

	t.Run("limit more than file size", func(t *testing.T) {
		srcFilePath := "testdata/input.txt"
		dstFilePath := "/tmp/output.txt"

		srcFileInfo, err := os.Stat(srcFilePath)
		require.NoError(t, err, "can't open src file")
		err = Copy(srcFilePath, dstFilePath, 0, srcFileInfo.Size()+100)
		require.NoError(t, err, "can't copy successfully")
		defer os.Remove(dstFilePath)

		dstFileInfo, err := os.Stat(dstFilePath)
		require.NoError(t, err, "can't get dst file info")
		require.Equal(t, srcFileInfo.Size(), dstFileInfo.Size(), "size of src and dst files must be equal")
	})

	t.Run("success copy", func(t *testing.T) {
		cases := []struct {
			name        string
			srcFilePath string
			dstFilePath string
			offset      int64
			limit       int64
			cmpFilePath string
		}{
			{
				name:        "copy with offset 0 and limit 0",
				srcFilePath: "testdata/input.txt",
				dstFilePath: "/tmp/out_offset0_limit0.txt",
				offset:      0,
				limit:       0,
				cmpFilePath: "testdata/out_offset0_limit0.txt",
			},
			{
				name:        "copy with offset 0 and limit 10",
				srcFilePath: "testdata/input.txt",
				dstFilePath: "/tmp/out_offset0_limit10.txt",
				offset:      0,
				limit:       10,
				cmpFilePath: "testdata/out_offset0_limit10.txt",
			},
			{
				name:        "copy with offset 0 and limit 1000",
				srcFilePath: "testdata/input.txt",
				dstFilePath: "/tmp/out_offset0_limit1000.txt",
				offset:      0,
				limit:       1000,
				cmpFilePath: "testdata/out_offset0_limit1000.txt",
			},
			{
				name:        "copy with offset 0 and limit 10000",
				srcFilePath: "testdata/input.txt",
				dstFilePath: "/tmp/out_offset0_limit10000.txt",
				offset:      0,
				limit:       10000,
				cmpFilePath: "testdata/out_offset0_limit10000.txt",
			},
			{
				name:        "copy with offset 100 and limit 1000",
				srcFilePath: "testdata/input.txt",
				dstFilePath: "/tmp/out_offset100_limit1000.txt",
				offset:      100,
				limit:       1000,
				cmpFilePath: "testdata/out_offset100_limit1000.txt",
			},
			{
				name:        "copy with offset 6000 and limit 1000",
				srcFilePath: "testdata/input.txt",
				dstFilePath: "/tmp/out_offset6000_limit1000.txt",
				offset:      6000,
				limit:       1000,
				cmpFilePath: "testdata/out_offset6000_limit1000.txt",
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := Copy(c.srcFilePath, c.dstFilePath, c.offset, c.limit)
				require.NoError(t, err, "can't copy successfully")
				defer os.Remove(c.dstFilePath)

				dstFileInfo, err := os.Stat(c.dstFilePath)
				require.NoError(t, err, "can't get dst file info")

				cmpFileInfo, err := os.Stat(c.cmpFilePath)
				require.NoError(t, err, "can't get cmp file info")

				require.Equal(t, cmpFileInfo.Size(), dstFileInfo.Size(), "size of cmp and dst files must be equal")
			})
		}
	})
}
