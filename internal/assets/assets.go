// package assets contains
//
// It must not be imported by the satelllite

package assets

import _ "embed"

//go:embed satellite-amd64-linux
var SatelliteAmd64Linux []byte
