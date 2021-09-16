package argus

import (
	"github.com/relvacode/argus/binary"
	"github.com/relvacode/argus/sensors"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func readBin(t *testing.T, name string) []byte {
	f, err := os.Open(filepath.Join("./_test", name))
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestSample_read(t *testing.T) {
	dump := readBin(t, "dump.bin")

	var sample Sample
	sample.read(binary.New(dump))

	assert.Equal(t, sample.Signature, Active)
	assert.Equal(t, sample.Version.Major, uint8(6))
	assert.Equal(t, sample.Version.MinorA, uint8(0))
	assert.Equal(t, sample.Version.MinorB, uint8(1))
	assert.Equal(t, sample.Version.Build, uint32(2507))

	assert.Equal(t, sample.StructureVersion, uint32(1))
	assert.Equal(t, sample.TotalMeasurementCount, uint32(74))
	assert.Equal(t, len(sample.Data), 74)

	assert.Equal(t, len(sample.Measurements(sensors.GPUTemperature)), 3)
	assert.Equal(t, sample.Measurements(sensors.GPUName)[0].Label, "NVIDIA NVIDIA GeForce RTX 3090")

	assert.Equal(t, len(sample.Measurements(0xFFFF)), 0)
}
