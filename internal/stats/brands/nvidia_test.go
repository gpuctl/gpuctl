package brands

import (
    "os"
    "reflect"
    "strings"
    "io/ioutil"
    "path/filepath"
    "testing"
    "encoding/json"
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
        res, err := ParseNvidiaSmi(dump)
        if err != nil {
            t.Errorf("Could not parse the nvidia-smi dump at %s: %v", fileloc, err)
            continue
        }

        r, err := res.FilterStatus()
        if err != nil {
            t.Errorf("Could not produce filtered status packet from nvidia-smi dump: %v (file %s) %v", err, fileloc, r)
            continue
        }

        j, err := json.Marshal(r)
        if err != nil {
            t.Errorf("Could not marshal status packet to JSON: %v (file %s)", err, fileloc)
            continue
        }
        // Read expected json
        var expected []byte

        {
            sp := strings.Split(filename, ".")
            resloc := testDataRoot + "/" + sp[0] + ".json"
            dump, err := ioutil.ReadFile(resloc)
            if err != nil {
                t.Fatalf("Could not read test result data at %s: %v", resloc, err)
            }
            expected = dump
        }

        j = append(j, 10) // HACK: annoying newline at the end of stored data...
        if !reflect.DeepEqual(expected, j) {
            t.Errorf("Parsed data did not match expected output (file %s): %v != %v", fileloc, j, expected)
        }

    }
}
