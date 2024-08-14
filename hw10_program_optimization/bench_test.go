package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	file, err := os.Open("testdata/users.dat.zip")
	if err != nil {
		b.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		b.Fatalf("failed to get file info: %v", err)
	}

	zipReader, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		b.Fatalf("failed to create zip reader: %v", err)
	}

	if len(zipReader.File) != 1 {
		b.Fatalf("expected 1 file in zip, got %d", len(zipReader.File))
	}

	zipFile := zipReader.File[0]
	rc, err := zipFile.Open()
	if err != nil {
		b.Fatalf("failed to open zip file: %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		b.Fatalf("failed to read data: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(data)
		_, err := GetDomainStat(r, "com")
		if err != nil {
			b.Fatalf("GetDomainStat failed: %v", err)
		}
	}
}
