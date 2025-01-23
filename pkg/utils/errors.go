package utils

import "errors"

var ErrUnsupportedVendor = errors.New("unsupported vendor")
var ErrUnsafeCommand = errors.New("Unsafe command detected")
var ErrEmptyCommand = errors.New("Empty command")
var ErrNoNvidiaDevice = errors.New("no nvidia device found")
