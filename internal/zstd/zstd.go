// Package zstd compresses data with the zstd C library. It is the example's
// cgo package: the library comes from nixpkgs, through the packageOverrides
// entry for this package in flake.nix.
package zstd

/*
#cgo pkg-config: libzstd
#include <zstd.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Version returns the version of the zstd library the program is linked
// against.
func Version() string {
	return C.GoString(C.ZSTD_versionString())
}

// Compress returns src as one zstd frame, at the library's default level.
func Compress(src []byte) ([]byte, error) {
	// The bound is positive even for empty input, so dst[0] exists.
	dst := make([]byte, C.ZSTD_compressBound(C.size_t(len(src))))
	var in unsafe.Pointer
	if len(src) > 0 {
		in = unsafe.Pointer(&src[0])
	}
	n := C.ZSTD_compress(unsafe.Pointer(&dst[0]), C.size_t(len(dst)), in, C.size_t(len(src)), C.ZSTD_CLEVEL_DEFAULT)
	if C.ZSTD_isError(n) != 0 {
		return nil, fmt.Errorf("zstd: %s", C.GoString(C.ZSTD_getErrorName(n)))
	}
	return dst[:n], nil
}
