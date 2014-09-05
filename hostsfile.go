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
	records []Record
}

// A single line in the hosts file
type Record struct {
	IpAddress net.IP
	Hostnames map[string]bool
	comment   string
	isBlank   bool
}

// Decodes the raw text of a hostsfile into a Hostsfile struct.
// Interface example from the image package.
func Decode(rdr io.Reader) (Hostsfile, error) {
	var h Hostsfile
	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)
		var r Record
		if len(line) == 0 {
			r.isBlank = true
		} else if line[0] == '#' {
			// comment line or blank line: skip it.
			r.comment = line
		} else {
			vals := strings.SplitN(line, " ", 2)
			if len(vals) <= 1 {
				return Hostsfile{}, fmt.Errorf("Invalid hostsfile entry: %s", line)
			}
			ip := net.ParseIP(vals[0])
			if ip == nil {
				return Hostsfile{}, fmt.Errorf("Invalid IP address: %s", vals[0])
			}
			r := Record{
				IpAddress: ip,
				Hostnames: map[string]bool{},
			}
			names := strings.Split(vals[1], " ")
			for i := 0; i < len(names); i++ {
				name := names[i]
				r.Hostnames[name] = true
			}
			h.records = append(h.records, r)
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
func (h *Hostsfile) Set(ip net.IP, hostname string) error {
	if ip == nil {
		return fmt.Errorf("Invalid IP address")
	}
	if len(hostname) == 0 {
		return fmt.Errorf("Hostname cannot be empty")
	}
	addKey := true
	for i := 0; i < len(h.records); i++ {
		record := h.records[i]
		if _, ok := record.Hostnames[hostname]; ok {
			if record.IpAddress.Equal(ip) {
				// tried to set a key that exists, nothing to do
				addKey = false
			} else {
				// delete the key and be sure to add a new record.
				delete(record.Hostnames, hostname)
				if len(record.Hostnames) == 0 {
					// delete the record
					h.records = append(h.records[:i], h.records[i+1:]...)
				}
				addKey = true
			}
		}
	}

	if addKey {
		nr := Record{
			IpAddress: ip,
			Hostnames: map[string]bool{hostname: true},
		}
		h.records = append(h.records, nr)
	}
	return nil
}

// Removes a hostname from the list. If the hostname is an alias,
func (h *Hostsfile) Remove(hostname string) error {
	return nil
}

// Return the text representation of the hosts file.
func Encode(w io.Writer, h Hostsfile) error {
	for _, record := range h.records {
		var toWrite string
		if len(record.comment) > 0 {
			toWrite = record.comment
		} else if record.isBlank {
			toWrite = ""
		} else {
			out := make([]string, len(record.Hostnames)+1)
			out[0] = record.IpAddress.String()
			i := 1
			for name, _ := range record.Hostnames {
				out[i] = name
				i++
			}
			toWrite = strings.Join(out, " ")
		}
		toWrite += "\n"
		_, err := w.Write([]byte(toWrite))
		if err != nil {
			return err
		}
	}
	return nil
}
