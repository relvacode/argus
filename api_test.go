package argus

import (
	"fmt"
	"github.com/relvacode/argus/sensors"
	"time"
)

func ExampleApi_Watch() {
	api, err := Open()
	if err != nil {
		panic(err)
	}

	defer api.Close()

	w := api.Watch(time.Millisecond * 100)
	defer w.Done()

	for data := range w.C {
		for _, m := range data.Measurements(sensors.GPUTemperature) {
			fmt.Println(m)
		}
	}

	err = w.Err()
	if err != nil {
		panic(err)
	}
}

func ExampleApi_Cached() {
	api, err := Open()
	if err != nil {
		panic(err)
	}

	defer api.Close()

	cached := api.Cached()

	for range time.NewTicker(time.Second).C {
		data, err := cached.Read()
		if err != nil {
			panic(err)
		}

		for _, m := range data.Measurements(sensors.GPUTemperature) {
			fmt.Println(m)
		}
	}
}
