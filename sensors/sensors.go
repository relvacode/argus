package sensors

//go:generate go run github.com/abice/go-enum -f $GOFILE --noprefix

// Type is an enumeration of all possible sensor types
/*
ENUM(
Invalid
Temperature // Temperatures of mainboard sensors, external fan controllers and AIOs
SyntheticTemperature // User defined synthetic temperature (mean, max, average, difference, etc..)
FanSpeedRPM // Fan speed of fans attached to mainboard channels, AIOs, external fan controllers and also pump speeds of AIOs
FanControlValue // If any fan or pump is controlled by Argus Monitor then the control value can be read from this
NetworkSpeed // Up/down speeds of network adapters if selected to be monitored inside Argus Monitor
CPUTemperature // The normal CPU temperature readings per core for Intel and the only one available for AMD
CPUTemperatureAdditional // Additional temperatures provided by the CPU like CCDx temperatures of AMD CPUs
CPUMultiplier // Multiplier value for every core
CPUFrequencyFSB // Core frequencies can be calculated by multiplying FSB frequency by the multipliers
GPUTemperature
GPUName // The name of the GPU (e.g. "Nvidia RTX3080")
GPULoad
GPUCoreClock
GPUMemoryClock
GPUShaderClock
GPUFanSpeedPercent
GPUFanSpeedRPM
GPUMemoryUsedPercent
GPUMemoryUsedMB
GPUPower
DiskTemperature
DiskTransferRate
CPULoad
RAMUsage
Battery
Length // The number of valid sensors. Is not a sensor itself
)
*/
type Type uint32

// Valid returns true if Type is a valid sensor type
func (x Type) Valid() bool {
	return x > Invalid && x < Length
}
