package brands

import (
    "os"
    "io/ioutil"
    "path/filepath"
    "testing"
)

const (
    testDataRoot  = "_testdata/nvidia"
    dataExtension = ".xml"
)

func TestNvidiaSmiXMLParsing(t *testing.T) {
    files, err := os.ReadDir(testDataRoot)
    if err != nil {
        t.Fatalf("Could not read test data root: %v", err)
    }
    for _, file := range files {
        filename := file.Name()
        if filepath.Ext(filename) != dataExtension {
            continue
        }
        fileloc := testDataRoot + "/" + filename
        dump, err := ioutil.ReadFile(fileloc)
        if err != nil {
            t.Fatalf("Could not read test data: %v", err)
        }
        if _, err := ParseNvidiaSmi(dump); err != nil {
            t.Errorf("Could not parse the nvidia-smi dump at %s", fileloc)
        }
    }
}

func TestNvidiaSmiParsingExtractsCorrectInformation(t *testing.T) {
    t.Fail()
}
