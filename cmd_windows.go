package main

import "os"

// Default location of hosts-file on Windows.
var hostsFile = os.Getenv("SystemRoot") + "\\System32\\drivers\\etc\\hosts"
