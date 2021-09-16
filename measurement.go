package argus

import (
	"fmt"
	"github.com/relvacode/argus/binary"
	"github.com/relvacode/argus/sensors"
)

const (
	labelArraySize        = 64 * binary.Uint16
	unitArraySize         = 32 * binary.Uint16
	measurementBufferSize = binary.Uint32 + labelArraySize + unitArraySize + binary.Float64 + binary.Uint32 + binary.Uint32
)

type Measurement struct {
	Type        sensors.Type
	Label       string
	Unit        string
	Value       float64
	DataIndex   uint32
	SensorIndex uint32
}

func (m *Measurement) read(r *binary.Reader) {
	m.Type = sensors.Type(r.Uint32())

	pos := r.Pos()
	m.Label = r.Utf16String()
	r.Seek(pos + labelArraySize)

	pos = r.Pos()
	m.Unit = r.Utf16String()
	r.Seek(pos + unitArraySize)

	m.Value = r.Float64()
	m.DataIndex = r.Uint32()
	m.SensorIndex = r.Uint32()
}

func (m Measurement) String() string {
	return fmt.Sprintf("%s.%q %f%s", m.Type, m.Label, m.Value, m.Unit)
}
