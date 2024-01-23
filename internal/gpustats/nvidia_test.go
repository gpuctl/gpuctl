package gpustats

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testDataRoot            = "_testdata/nvidia"
	workingDataExtension    = ".xml"
	faultyCallDataExtension = ".faultycall"
	corruptedDataExtension  = ".corruptedxml"
)

func TestNvidiaSmiXMLParsing(t *testing.T) {
	files, err := os.ReadDir(testDataRoot)
	if err != nil {
		t.Fatalf("Could not read test data root: %v", err)
	}
	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) != workingDataExtension {
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

		j = append(j, '\n') // HACK: annoying newline at the end of stored data...
		if !bytes.Equal(expected, j) {
			t.Errorf("Parsed data did not match expected output (file %s):\n%s!=\n%s", fileloc, j, expected)
		}

	}
}

func TestNvidiaSmiFaultyInput(t *testing.T) {
	files, err := os.ReadDir(testDataRoot)
	if err != nil {
		t.Fatalf("Could not read test data root: %v", err)
	}

	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) != faultyCallDataExtension {
			continue
		}
		fileloc := testDataRoot + "/" + filename
		dump, err := ioutil.ReadFile(fileloc)
		if err != nil {
			t.Fatalf("Could not read test data: %v", err)
		}
		_, err = ParseNvidiaSmi(dump)
		if err == nil {
			t.Errorf("Accepted invalid nvidia dump (file %s)", fileloc)
			continue
		}
	}
}

func TestNvidiaSmiInvalidDataParse(t *testing.T) {
	files, err := os.ReadDir(testDataRoot)
	if err != nil {
		t.Fatalf("Could not read test data root: %v", err)
	}

	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) != ".corruptedxml" {
			continue
		}
		fileloc := testDataRoot + "/" + filename
		dump, err := ioutil.ReadFile(fileloc)
		if err != nil {
			t.Fatalf("Could not read test data: %v", err)
		}
		smi, err := ParseNvidiaSmi(dump)
		if err != nil {
			t.Errorf("Could not parse file %s: %v", fileloc, err)
			continue
		}

		_, err = smi.FilterStatus()
		if err == nil {
			t.Errorf("Accepted mangled data in parsing fields of nvidia-smi data (file %s)", fileloc)
			continue
		}
	}
}
