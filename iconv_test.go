package iconv

import "fmt"
import "testing"

var data = []struct{original, encoded, encoding string} {
	{"привет", "\xef\xf0\xe8\xe2\xe5\xf2", "Windows-1251"},
	{"привет", "\x04\x3f\x04\x40\x04\x38\x04\x32\x04\x35\x04\x42", "UTF-16BE"},
	{"a", "\x00\x61", "UTF-16BE"},
	{"これは漢字です。", "\x82\xb1\x82\xea\x82\xcd\x8a\xbf\x8e\x9a\x82\xc5\x82\xb7\x81B", "SJIS"},
	{"これは漢字です。", "S0\x8c0o0\"oW[g0Y0\x020", "UTF-16LE"},
	{"これは漢字です。", "0S0\x8c0oo\"[W0g0Y0\x02", "UTF-16BE"},
	{"€1 is cheap", "\xa41 is cheap", "ISO-8859-15"},
	{"", "", "SJIS"}}

func TestIconv(t *testing.T) {
	for _, item := range data {
		c := Open(item.encoding, "UTF-8")
		s := c.Conv(item.original)
		if s != item.encoded { t.Error("") }
		c.Close()
	}

	for _, item := range data {
		c := Open("UTF-8", item.encoding)
		s := c.Conv(item.encoded)
		if s != item.original { t.Error("") }
		c.Close()
	}
}

func TestIconvMultiple(t *testing.T) {
	item := data[0]
	c := Open(item.encoding, "UTF-8"); defer c.Close()

	s := c.Conv(item.original)
	if s != item.encoded { t.Error("") }

	s = c.Conv(item.original)
	if s != item.encoded { t.Error("") }
}

func brake(c *Iconv) {
	defer func() {
		e := recover()
		fmt.Println("brake: recover: e =", e)
	}()

	broken := "Р С—РЎР‚Р С‘Р Р†Р ВµРЎвЂљ\x0a"
	c.Conv(broken)
}

// FIXME utf-16be is not stateful, it seems. So test does not work
func TestAfterBroken(t *testing.T) {
	valid := "\x04\x3f\x04\x40\x04\x38\x04\x32\x04\x35\x04\x42"

	c := Open("UTF-8", "UTF-16BE"); defer c.Close()
	brake(c)

	s := c.Conv(valid)
	if s != "привет" { t.Error("") }
}
