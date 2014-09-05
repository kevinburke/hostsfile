package hostsfile

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

// Represents a hosts file. Records match a single line in the file.
type Hostsfile struct {
	Records []Record
	raw     []string
}

// A single line in the hosts file
type Record struct {
	IpAddress net.IP
	Hostname  string
	Aliases   []string
}

// Decodes the raw text of a hostsfile into a Hostsfile struct.
// Interface example from the image package.
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

// Adds a record to the list. If the hostname is present with a different IP
// address, it will be reassigned. If the record is already present with the
// same hostname/IP address data, it will not be added again.
func (h *Hostsfile) Set(r Record) {

}

// Removes a hostname from the list. If the hostname is an alias,
func (h *Hostsfile) Remove(hostname string) error {

}

// Writes a hostsfile to a string
func Encode(w io.Writer, h Hostsfile) error {

}
