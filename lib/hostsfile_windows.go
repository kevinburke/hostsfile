package hostsfile

import "os"

// OS-specific default hosts-file location.
var Location = os.Getenv("SystemRoot") + "\\System32\\drivers\\etc\\hosts"

// OS-specific newline character(s).
const eol = "\r\n"
