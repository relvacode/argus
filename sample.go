package argus

import (
	"fmt"
	"github.com/relvacode/argus/binary"
	"github.com/relvacode/argus/sensors"
)

const (
	sensorDataArraySize = 512
	versionBufferSize   = (binary.Uint8 * 4) + binary.Uint32
	sampleBufferSize    = binary.Uint32 + versionBufferSize + binary.Uint32 + binary.Uint32 + (binary.Uint32 * sensors.Length) + (binary.Uint32 * sensors.Length) + binary.Uint32 + (sensorDataArraySize * measurementBufferSize)
	cycleCounterOffset  = binary.Uint32 + versionBufferSize + binary.Uint32
)

const (
	// Active is the value of Sample.Signature if Argus Monitor is active
	Active   uint32 = 0x4D677241
	Inactive uint32 = 0x00000000
)

type Version struct {
	Major  uint8
	MinorA uint8
	MinorB uint8
	Extra  uint8
	Build  uint32
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d.%d.%d", v.Major, v.MinorA, v.MinorB, v.Extra, v.Build)
}

// Sample is a snapshot of measurements taken from Argus Monitor.
type Sample struct {
	Signature uint32
	Version
	StructureVersion uint32
	CycleCounter     uint32

	// OffsetForSensorType describes the offsets for each sensors.Type.
	OffsetForSensorType [sensors.Length]uint32

	// CountsForSensorType describes the length of measurements for each sensors.Type.
	CountsForSensorType [sensors.Length]uint32

	TotalMeasurementCount uint32
	Data                  []Measurement
}

func (sample *Sample) read(r *binary.Reader) {
	sample.Signature = r.Uint32()

	sample.Version.Major = r.Uint8()
	sample.Version.MinorA = r.Uint8()
	sample.Version.MinorB = r.Uint8()
	sample.Version.Extra = r.Uint8()
	sample.Version.Build = r.Uint32()

	sample.StructureVersion = r.Uint32()
	sample.CycleCounter = r.Uint32()

	for i := 0; i < len(sample.OffsetForSensorType); i++ {
		sample.OffsetForSensorType[i] = r.Uint32()
	}

	for i := 0; i < len(sample.CountsForSensorType); i++ {
		sample.CountsForSensorType[i] = r.Uint32()
	}

	sample.TotalMeasurementCount = r.Uint32()

	sample.Data = make([]Measurement, 0, sensorDataArraySize)
	for i := 0; i < sensorDataArraySize && i < int(sample.TotalMeasurementCount); i++ {
		var measurement Measurement
		(&measurement).read(r)
		sample.Data = append(sample.Data, measurement)
	}
}

// Measurements returns all measurements for a given sensors.Type.
func (sample *Sample) Measurements(typ sensors.Type) []Measurement {
	if typ >= sensors.Length {
		return nil
	}

	var (
		dataOffset  = int(sample.OffsetForSensorType[typ])
		sensorCount = int(sample.CountsForSensorType[typ])
	)

	return sample.Data[dataOffset : dataOffset+sensorCount]
}
