# go-hostsfile

This library will help you manipulate your /etc/hosts file. A description of
the API [can be found at godoc][godoc].

## Command Line Usage

You can install a command line binary:

```
go get github.com/kevinburke/hostsfile
```

and use it like so:

```
hostsfile add www.facebook.com www.twitter.com www.adroll.com 127.0.0.1
hostsfile remove www.facebook.com www.twitter.com www.adroll.com
```

You may need to run the above commands as root to write to `/etc/hosts` (which
is modified atomically).

To print the new file to stdout, instead of writing it:

```
hostsfile add --dry-run www.facebook.com www.twitter.com www.adroll.com 127.0.0.1
```

You can also pipe a hostsfile in:

```
cat /etc/hosts | hostsfile add --dry-run www.facebook.com www.twitter.com www.adroll.com 127.0.0.1
```

Or specify a file to read from at the command line:

```
hostsfile add --file=sample-hostsfile www.facebook.com www.twitter.com www.adroll.com 127.0.0.1
```

## Library Usage

You can also call the functions in this library from Go code.

```go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	hostsfile "github.com/kevinburke/hostsfile/lib"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	f, err := os.Open("/etc/hosts")
	checkError(err)
	h, err := hostsfile.Decode(f)
	checkError(err)

	local, err := net.ResolveIPAddr("ip", "127.0.0.1")
	checkError(err)
	// Necessary for sites like facebook & gmail that resolve ipv6 addresses,
	// if your network supports ipv6
	ip6, err := net.ResolveIPAddr("ip", "::1")
	checkError(err)
	h.Set(*local, "www.facebook.com")
	h.Set(*ip6, "www.facebook.com")
	h.Set(*local, "news.ycombinator.com")
	h.Set(*ip6, "news.ycombinator.com")


	// Write to a temporary file and then atomically copy it into place.
	tmpf, err := ioutil.TempFile("/tmp", "hostsfile-temp")
	checkError(err)

	err = hostsfile.Encode(tmpf, h)
	checkError(err)

	err = os.Chmod(tmp.Name(), 0644)
	checkError(err)

	err = os.Rename(tmp.Name(), "/etc/hosts")
	checkError(err)
	fmt.Println("done")
}
```

[godoc]: https://godoc.org/github.com/kevinburke/hostsfile
