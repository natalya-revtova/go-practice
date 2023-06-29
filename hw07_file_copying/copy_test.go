package main

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	description     string
	sourceFile      string
	destinationFile string
	offset          int64
	limit           int64
	want            string
}

func TestCopy_Success(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				"offset 0, limit 0",
				"testdata/input.txt",
				"out.txt",
				0,
				0,
				"testdata/out_offset0_limit0.txt",
			},
			{
				"offset 10, limit 0",
				"testdata/input.txt",
				"/tmp/out.txt",
				0,
				10,
				"testdata/out_offset0_limit10.txt",
			},
			{
				"offset 1000, limit 0",
				"testdata/input.txt",
				"/tmp/out.txt",
				0,
				1000,
				"testdata/out_offset0_limit1000.txt",
			},
			{
				"offset 10000, limit 0",
				"testdata/input.txt",
				"/tmp/out.txt",
				0,
				10000,
				"testdata/out_offset0_limit10000.txt",
			},
			{
				"offset 100, limit 1000",
				"testdata/input.txt",
				"/tmp/out.txt",
				100,
				1000,
				"testdata/out_offset100_limit1000.txt",
			},
			{
				"offset 6000, limit 1000",
				"testdata/input.txt",
				"/tmp/out.txt",
				6000,
				1000,
				"testdata/out_offset6000_limit1000.txt",
			},
		}

		for _, tt := range tests {
			t.Run(tt.description, func(t *testing.T) {
				err := Copy(tt.sourceFile, tt.destinationFile, tt.offset, tt.limit)
				require.NoError(t, err)

				hashOutResult := getHash(t, tt.destinationFile)
				hashOutTest := getHash(t, tt.want)
				assert.Equal(t, hashOutResult, hashOutTest)
			})
		}
	})

	t.Run("invalid files", func(t *testing.T) {
		srcAbs, err := filepath.Abs("testdata/input.txt")
		require.NoError(t, err)

		tests := []testCase{
			{
				"empty source file parameter",
				"",
				"/tmp/out.txt",
				0,
				0,
				"",
			},
			{
				"empty destination file parameter",
				"testdata/input.txt",
				"",
				0,
				0,
				"",
			},
			{
				"equal source file and destination files parameters",
				srcAbs,
				"testdata/input.txt",
				0,
				0,
				"",
			},
			{
				"invalid source file",
				"testdata/",
				"/tmp/out.txt",
				0,
				0,
				"",
			},
			{
				"unknown source file length",
				"/dev/urandom",
				"/tmp/out.txt",
				0,
				0,
				"",
			},
		}

		for _, tt := range tests {
			t.Run(tt.description, func(t *testing.T) {
				err := Copy(tt.sourceFile, tt.destinationFile, tt.offset, tt.limit)
				assert.ErrorIs(t, err, ErrUnsupportedFile)
			})
		}
	})

	t.Run("invalid offset", func(t *testing.T) {
		tests := []testCase{
			{
				"offset is more than source file length",
				"testdata/input.txt",
				"/tmp/out.txt",
				10000,
				0,
				"",
			},
		}

		for _, tt := range tests {
			t.Run(tt.description, func(t *testing.T) {
				err := Copy(tt.sourceFile, tt.destinationFile, tt.offset, tt.limit)
				assert.ErrorIs(t, err, ErrOffsetExceedsFileSize)
			})
		}
	})
}

func getHash(t *testing.T, path string) []byte {
	t.Helper()

	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	require.NoError(t, err)

	return hash.Sum(nil)
}
