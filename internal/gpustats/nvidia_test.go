package gpustats

import (
	//	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"strings"
	"testing"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

const (
	testDataRoot            = "_testdata/nvidia"
	workingDataExtension    = ".xml"
	faultyCallDataExtension = ".faultycall"
	corruptedDataExtension  = ".corruptedxml"
)

func TestNvidiaSmiXMLParsing(t *testing.T) {
	t.Parallel()
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
		dump, err := os.ReadFile(fileloc)
		if err != nil {
			t.Fatalf("Could not read test data: %v", err)
		}
		res, err := ParseNvidiaSmi(dump)
		if err != nil {
			t.Errorf("Could not parse the nvidia-smi dump at %s: %v", fileloc, err)
			continue
		}

		stats, err := res.ExtractGPUStatSample()
		if err != nil {
			t.Errorf("Could not extract GPU status from nvidia-smi dump: %v (file %s)", err, fileloc)
			continue
		}

		info, err := res.ExtractGPUInfo()
		if err != nil {
			t.Errorf("Could not extract general GPU info from nvidia-smi dump: %v (file %s)", err, fileloc)
			continue
		}

		result := uplink.GpuStatsUpload{Hostname: "", GPUInfos: info, Stats: stats}
		resultJson, err := json.Marshal(result)

		if err != nil {
			t.Errorf("Could not marshal status packet to JSON: %v (file %s)", err, fileloc)
			continue
		}

		// Compare parsed resultJson data with expected output
		sp := strings.Split(filename, ".")
		resloc := testDataRoot + "/" + sp[0] + ".json"
		expected_dump, err := os.ReadFile(resloc)
		if err != nil {
			t.Fatalf("Could not read test result data at %s: %v", resloc, err)
		}
		var expected uplink.GpuStatsUpload
		err = json.Unmarshal(expected_dump, &expected)
		if err != nil {
			t.Fatalf("Could not unmarshal test result data at %s: %v", resloc, err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Result data did not match expected. \nGot      %s \nexpected %s \n(file: %s)", resultJson, string(expected_dump), fileloc)
		}

		/*
			// Read expected json
			var expected []byte

			{
				sp := strings.Split(filename, ".")
				resloc := testDataRoot + "/" + sp[0] + ".json"
				expected = dump
			}

			j = append(j, '\n') // HACK: annoying newline at the end of stored data...
			if !bytes.Equal(expected, j) {
				t.Errorf("Parsed data did not match expected output (file %s):\n%s!=\n%s", fileloc, j, expected)
			}
		*/

	}
}

func TestNvidiaSmiFaultyInput(t *testing.T) {
	t.Parallel()
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
		dump, err := os.ReadFile(fileloc)
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
	t.Parallel()
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
		dump, err := os.ReadFile(fileloc)
		if err != nil {
			t.Fatalf("Could not read test data: %v", err)
		}
		smi, err := ParseNvidiaSmi(dump)
		if err != nil {
			t.Errorf("Could not parse file %s: %v", fileloc, err)
			continue
		}

		_, err = smi.ExtractGPUStatSample()
		if err == nil {
			t.Errorf("Accepted mangled data in parsing fields of nvidia-smi data (file %s)", fileloc)
			continue
		}
	}
}
