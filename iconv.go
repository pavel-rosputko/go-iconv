package iconv

import (
	"bytes"
	"os"
)

// #include <iconv.h>
// #include <stdlib.h>
// #include <errno.h>
import "C"

import (
	"unsafe"
)

type Iconv struct {
	p	C.iconv_t
}

func Open(to, from string) *Iconv {
	cto := C.CString(to); defer C.free(unsafe.Pointer(cto))
	cfrom := C.CString(from); defer C.free(unsafe.Pointer(cfrom))
	cc, e := C.iconv_open(cto, cfrom)
	if int(uintptr(cc)) == -1 { panic(e) }
	return &Iconv{cc}
}

func (c *Iconv) Close() {
	r, e := C.iconv_close(c.p)
	if r == -1 { panic(e) }
}

// FIXME loop forever if inLen not decreasing
// FIXME if error reset state ?
func (c *Iconv) Conv(s string) string {
	if len(s) == 0 { return "" }

	inBytes := []byte(s)

	bufferLen := len(inBytes)
	if bufferLen < 64 { bufferLen = 64 } // XXX one character never exceeds this ?
	buffer := make([]byte, bufferLen)

	outBuffer := bytes.NewBuffer(nil)

	inBytesPtr := &inBytes[0]
	inLen := C.size_t(len(inBytes))
	for inLen > 0 {
		bufferPtr := &buffer[0]
		bufferLen := C.size_t(len(buffer))
		r, e := C.iconv(c.p,
			(**C.char)(unsafe.Pointer(&inBytesPtr)), &inLen,
			(**C.char)(unsafe.Pointer(&bufferPtr)), &bufferLen)
		if int(r) == -1 && e != os.Errno(int(C.E2BIG)) { panic(e) }

		outBuffer.Write(buffer[:len(buffer) - int(bufferLen)])
	}

	// XXX if input is incomplete we need to reset iconv state (if it has)
	// or call a version that returns shift sequence ?
	r, e := C.iconv(c.p,
		(**C.char)(nil), (*C.size_t)(nil),
		(**C.char)(nil), (*C.size_t)(nil))
	if int(r) == -1 { panic(e) }

	return outBuffer.String()
}
