package hostsfile

import (
	"bytes"
	"fmt"
	"net"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\texp: %#v\n\tgot: %#v\033[39m\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()
	sampledata := "127.0.0.1 foobar\n# this is a comment\n10.0.0.1 anotheralias"
	h, err := Decode(strings.NewReader(sampledata))
	if err != nil {
		t.Error(err.Error())
	}
	firstRecord := h.records[0]

	equals(t, firstRecord.IpAddress.IP.String(), "127.0.0.1")
	equals(t, firstRecord.Hostnames["foobar"], true)
	equals(t, len(firstRecord.Hostnames), 1)

	aliasSample := "127.0.0.1 name alias1 alias2 alias3"
	h, err = Decode(strings.NewReader(aliasSample))
	ok(t, err)
	hns := h.records[0].Hostnames
	equals(t, len(hns), 4)
	equals(t, hns["alias3"], true)

	badline := strings.NewReader("blah")
	h, err = Decode(badline)
	if err == nil {
		t.Error("expected Decode(\"blah\") to return invalid, got no error")
	}
	if err.Error() != "Invalid hostsfile entry: blah" {
		t.Errorf("expected Decode(\"blah\") to return invalid, got %s", err.Error())
	}

	h, err = Decode(strings.NewReader("##\n127.0.0.1\tlocalhost    2nd-alias"))
	ok(t, err)
	equals(t, h.records[0].Hostnames["2nd-alias"], true)

	h, err = Decode(strings.NewReader("##\n127.0.0.1\tlocalhost # a comment"))
	ok(t, err)
	equals(t, h.records[0].Hostnames["#"], false)
	equals(t, h.records[0].Hostnames["a"], false)
}

func sample(t *testing.T) Hostsfile {
	one27, err := net.ResolveIPAddr("ip", "127.0.0.1")
	ok(t, err)
	one92, err := net.ResolveIPAddr("ip", "192.168.0.1")
	ok(t, err)
	return Hostsfile{
		records: []Record{
			Record{
				IpAddress: *one27,
				Hostnames: map[string]bool{"foobar": true},
			},
			Record{
				IpAddress: *one92,
				Hostnames: map[string]bool{"bazbaz": true, "blahbar": true},
			},
		},
	}
}

func comment(t *testing.T) Hostsfile {
	one92, err := net.ResolveIPAddr("ip", "192.168.0.1")
	ok(t, err)
	return Hostsfile{
		records: []Record{
			Record{
				comment: "# Don't delete this line!",
			},
			Record{
				comment: "shouldnt matter",
				isBlank: true,
			},
			Record{
				IpAddress: *one92,
				Hostnames: map[string]bool{"bazbaz": true},
			},
		},
	}
}

func TestEncode(t *testing.T) {
	t.Parallel()
	b := new(bytes.Buffer)
	err := Encode(b, sample(t))
	ok(t, err)
	equals(t, b.String(), "127.0.0.1 foobar\n192.168.0.1 bazbaz blahbar\n")

	b.Reset()
	err = Encode(b, comment(t))
	ok(t, err)
	equals(t, b.String(), "# Don't delete this line!\n\n192.168.0.1 bazbaz\n")
}

func TestRemove(t *testing.T) {
	t.Parallel()
	hCopy := sample(t)
	equals(t, len(hCopy.records[1].Hostnames), 2)
	hCopy.Remove("bazbaz")
	equals(t, len(hCopy.records[1].Hostnames), 1)
	ok := hCopy.records[1].Hostnames["blahbar"]
	assert(t, ok, "item \"blahbar\" not found in %v", hCopy.records[1].Hostnames)
	hCopy.Remove("blahbar")
	equals(t, len(hCopy.records), 1)
}

func TestSet(t *testing.T) {
	t.Parallel()
	hCopy := sample(t)
	one0, err := net.ResolveIPAddr("ip", "10.0.0.1")
	ok(t, err)
	hCopy.Set(*one0, "tendot")
	equals(t, len(hCopy.records), 3)
	equals(t, hCopy.records[2].Hostnames["tendot"], true)
	equals(t, hCopy.records[2].IpAddress.String(), "10.0.0.1")

	// appending same element shouldn't change anything
	hCopy.Set(*one0, "tendot")
	equals(t, len(hCopy.records), 3)

	one92, err := net.ResolveIPAddr("ip", "192.168.3.7")
	ok(t, err)
	hCopy.Set(*one92, "tendot")
	equals(t, hCopy.records[2].IpAddress.String(), "192.168.3.7")
}
