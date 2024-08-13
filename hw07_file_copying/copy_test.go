package main

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name        string
		srcFile     string
		dstFile     string
		offset      int64
		limit       int64
		expectError bool
	}{
		{
			name:        "Copy whole file",
			srcFile:     "testdata/input.txt",
			dstFile:     "/tmp/destination.txt",
			offset:      0,
			limit:       0,
			expectError: false,
		},
		{
			name:        "Copy with offset",
			srcFile:     "testdata/input.txt",
			dstFile:     "/tmp/destination_offset.txt",
			offset:      100,
			limit:       0,
			expectError: false,
		},
		{
			name:        "Copy with limit",
			srcFile:     "testdata/input.txt",
			dstFile:     "/tmp/destination_limit.txt",
			offset:      0,
			limit:       500,
			expectError: false,
		},
		{
			name:        "Copy with offset and limit",
			srcFile:     "testdata/input.txt",
			dstFile:     "/tmp/destination_offset_limit.txt",
			offset:      100,
			limit:       500,
			expectError: false,
		},
		{
			name:        "Offset exceeds file size",
			srcFile:     "testdata/input.txt",
			dstFile:     "/tmp/destination_offset_exceeds.txt",
			offset:      10000,
			limit:       0,
			expectError: true,
		},
		{
			name:        "Unsupported file",
			srcFile:     "/dev/urandom",
			dstFile:     "/tmp/destination_unsupported.txt",
			offset:      0,
			limit:       0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.srcFile, tt.dstFile, tt.offset, tt.limit)
			if (err != nil) != tt.expectError {
				t.Errorf("Copy() error = %v, expectError %v", err, tt.expectError)
			}
			if !tt.expectError {
				if _, err := os.Stat(tt.dstFile); os.IsNotExist(err) {
					t.Errorf("Destination file %s does not exist", tt.dstFile)
				} else {
					os.Remove(tt.dstFile)
				}
			}
		})
	}
}
