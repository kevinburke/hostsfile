# hostsfile

This library, and the associated command line binary, will help you manipulate
your /etc/hosts file. Both the library and the binary will leave comments
and other metadata in the /etc/hosts file as is, appending or removing only
the lines that you want changed. A description of the API [can be found at
godoc][godoc].

## Installation

On Mac, install via Homebrew:

```
brew install kevinburke/safe/hostsfile
```


If you have a Go development environment, you can install via source code:

    go get github.com/kevinburke/hostsfile@latest

## Command Line Usage

Easily add and remove entries from /etc/hosts.

```
# Assign 127.0.0.1 to all of the given hostnames
hostsfile add www.facebook.com www.twitter.com www.adroll.com 127.0.0.1
# Remove all hostnames from /etc/hosts
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

You can also call the functions in this library from Go code. Here's an example
where a hosts file is read, modified, and atomically written back to disk.

```go
package main

import (
	"bytes"
	"fmt"
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
	tmpf, err := os.CreateTemp("/tmp", "hostsfile-temp")
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
