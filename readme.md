# go-hostsfile

This library will help you manipulate your /etc/hosts file. A description of
the API [can be found at godoc][godoc].

## Sample Usage

```go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/kevinburke/hostsfile"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	file, err := os.Open("/etc/hosts")
	checkError(err)
	h, err := hostsfile.Decode(file)
	checkError(err)
	local, err := net.ResolveIPAddr("ip", "127.0.0.1")
	checkError(err)
	h.Set(*local, "www.facebook.com")
	h.Set(*local, "news.ycombinator.com")
	buf := new(bytes.Buffer)
	hostsfile.Encode(buf, h)

	// Write to a temporary file and then atomically copy it into place.
	tmp, err := ioutil.TempFile("/tmp", "hostsfile-temp")
	ioutil.WriteFile(tmp.Name(), buf.Bytes(), 0644)
	checkError(err)
	err = ioutil.WriteFile(tmp.Name(), buf.Bytes(), 0644)
	checkError(err)
	err = os.Rename(tmp.Name(), "/etc/hosts")
	checkError(err)
	fmt.Println("done")
}
```

[godoc]: http://godoc.org/github.com/kevinburke/hostsfile
