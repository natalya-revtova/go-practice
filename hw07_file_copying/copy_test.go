package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	src    string
	dest   string
	offset int64
	limit  int64
	want   string
}

func TestCopy_Success(t *testing.T) {
	tests := []testCase{
		{"testdata/input.txt", "/tmp/out.txt", 0, 0, "testdata/out_offset0_limit0.txt"},
		{"testdata/input.txt", "/tmp/out.txt", 0, 10, "testdata/out_offset0_limit10.txt"},
		{"testdata/input.txt", "/tmp/out.txt", 0, 1000, "testdata/out_offset0_limit1000.txt"},
		{"testdata/input.txt", "/tmp/out.txt", 0, 10000, "testdata/out_offset0_limit10000.txt"},
		{"testdata/input.txt", "/tmp/out.txt", 100, 1000, "testdata/out_offset100_limit1000.txt"},
		{"testdata/input.txt", "/tmp/out.txt", 6000, 1000, "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tt := range tests {
		err := Copy(tt.src, tt.dest, tt.offset, tt.limit)
		require.NoError(t, err)

		hashOutResult, _ := getHash(tt.dest)
		hashOutTest, _ := getHash(tt.want)
		require.Equal(t, hashOutResult, hashOutTest)

	}
}

func TestCopy_InvalidFiles(t *testing.T) {
	tests := []testCase{
		{"", "/tmp/out.txt", 0, 0, ""},
		{"testdata/input.txt", "", 0, 0, ""},
		{"testdata/input.txt", "testdata/input.txt", 0, 0, ""},
		{"testdata/", "/tmp/out.txt", 0, 0, ""},
		{"/dev/urandom", "/tmp/out.txt", 0, 0, ""},
	}

	for _, tt := range tests {
		err := Copy(tt.src, tt.dest, tt.offset, tt.limit)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	}
}

func TestCopy_InvalidOffset(t *testing.T) {
	tests := []testCase{
		{"testdata/input.txt", "/tmp/out.txt", 10000, 0, ""},
	}

	for _, tt := range tests {
		err := Copy(tt.src, tt.dest, tt.offset, tt.limit)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	}
}

func getHash(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to create hash: %w", err)
	}

	return hash.Sum(nil), nil
}
