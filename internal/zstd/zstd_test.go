package zstd

import (
	"bytes"
	"strings"
	"testing"
)

// frameMagic starts every zstd frame.
var frameMagic = []byte{0x28, 0xb5, 0x2f, 0xfd}

func TestVersion(t *testing.T) {
	if v := Version(); !strings.Contains(v, ".") {
		t.Fatalf("Version() = %q, want a dotted version", v)
	}
}

func TestCompress(t *testing.T) {
	src := bytes.Repeat([]byte("hello, zstd. "), 200)
	got, err := Compress(src)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(got, frameMagic) {
		t.Fatalf("output starts % x, want the zstd frame magic % x", got[:4], frameMagic)
	}
	if len(got) >= len(src) {
		t.Fatalf("compressed %d repetitive bytes to %d", len(src), len(got))
	}
}

func TestCompressEmpty(t *testing.T) {
	got, err := Compress(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(got, frameMagic) {
		t.Fatalf("output of empty input is % x, want a zstd frame", got)
	}
}
