package main

import (
	_"embed"
)

//go:embed certs/rootCA.pem
var rootCAPub []byte

//go:embed certs/private_key.pem
var rootCAPem []byte

//go:embed placeholder.png
var placeholder []byte

