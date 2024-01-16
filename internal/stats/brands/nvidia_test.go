package brands

import (
    "io/ioutil"
    "testing"
)

func TestNvidiaSmiXMLParsing(t *testing.T) {
    v, err := ioutil.ReadFile("filename") //read the content of file
    if err != nil {
        return
    }
	a := brands.ParseNvidiaSmi(v)
}
