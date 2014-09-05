package hostsfile

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

// Represents a hosts file
type Hostsfile struct {
	Records []Record
	raw     []string
}

type Record struct {
	IpAddress net.IP
	Hostname  string
	Aliases   []string
}

func Decode(r io.Reader) (Hostsfile, error) {
	var h Hostsfile
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)
		h.raw = append(h.raw, line)
		if len(line) == 0 || line[0] == '#' {
			// comment line or blank line: skip it.
		} else {
			vals := strings.SplitN(line, " ", 3)
			if len(vals) <= 1 {
				return Hostsfile{}, fmt.Errorf("Invalid hostsfile entry: %s", line)
			}
			ip := net.ParseIP(vals[0])
			if ip == nil {
				return Hostsfile{}, fmt.Errorf("Invalid IP address: %s", vals[0])
			}
			r := Record{
				IpAddress: ip,
				Hostname:  vals[1],
			}
			if len(vals) > 2 {
				r.Aliases = strings.Split(vals[2], " ")
			}
			h.Records = append(h.Records, r)
		}
	}
	if err := scanner.Err(); err != nil {
		return Hostsfile{}, err
	}
	return h, nil
}
