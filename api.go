// Package argus implements the Argus Monitor Data API interface.
// Reference C++ implementation https://github.com/argotronic/argus_data_api
package argus

import (
	"errors"
	"github.com/relvacode/argus/binary"
	"golang.org/x/sys/windows"
	"sync"
	"time"
	"unsafe"
)

const (
	mappingInterface      = "Global\\ARGUSMONITOR_DATA_INTERFACE"
	mappingInterfaceMutex = "Global\\ARGUSMONITOR_DATA_INTERFACE_MUTEX"
	mappingSize           = 1024 * 1024
)

func Open() (*Api, error) {
	mappingNameBuf, _ := windows.UTF16FromString(mappingInterface)
	dataHandle, err := syscallOpenFileMappingW(windows.FILE_MAP_READ|windows.FILE_MAP_WRITE, false, &mappingNameBuf[0])
	if err != nil {
		return nil, err
	}

	addr, err := windows.MapViewOfFile(dataHandle, windows.FILE_MAP_READ|windows.FILE_MAP_WRITE, 0, 0, mappingSize)
	if err != nil {
		_ = windows.CloseHandle(dataHandle)
		return nil, err
	}

	mutexNameBuf, _ := windows.UTF16FromString(mappingInterfaceMutex)
	mutexHandle, err := windows.OpenMutex(windows.READ_CONTROL|windows.MUTANT_QUERY_STATE|windows.SYNCHRONIZE, false, &mutexNameBuf[0])
	if err != nil {
		_ = windows.UnmapViewOfFile(addr)
		_ = windows.CloseHandle(dataHandle)
		return nil, errors.New("failed to acquire mutex")
	}

	return &Api{
		handle: dataHandle,
		mutex:  mutexHandle,
		addr:   addr,
	}, nil
}

type Api struct {
	handle windows.Handle
	mutex  windows.Handle
	addr   uintptr
}

// Close closes any open handles to the shared memory interface of Argus Monitor.
func (a *Api) Close() error {
	unmapError := windows.UnmapViewOfFile(a.addr)
	closeError := windows.CloseHandle(a.handle)

	if unmapError != nil {
		return unmapError
	}
	if closeError != nil {
		return closeError
	}

	return nil
}

// CycleCounter is fast method to check the CycleCounter for the current Sample in Argus Monitor.
func (a *Api) CycleCounter() uint32 {
	return *(*uint32)(unsafe.Pointer(a.addr + cycleCounterOffset))
}

// Read returns sensor data from the Argus Monitor interface.
// It always reads a new sample from Argus Monitor, even if the sample hasn't changed since the last Read call.
// Not safe for concurrent access.
// Use Api.Cached for faster access.
func (a *Api) Read() (*Sample, error) {
	_, err := windows.WaitForSingleObject(a.mutex, windows.INFINITE)
	if err != nil {
		return nil, err
	}

	var sample Sample
	var buffer = unsafe.Slice((*byte)(unsafe.Pointer(a.addr)), sampleBufferSize)
	sample.read(binary.New(buffer))

	err = windows.ReleaseMutex(a.mutex)
	if err != nil {
		return nil, err
	}

	return &sample, nil
}

type Watcher struct {
	// C is a channel that emits new samples after each change.
	// It will continue to emit changes until the Watcher is done or an error occurs.
	C <-chan *Sample

	done chan struct{}
	once sync.Once

	mx  sync.Mutex
	err error
}

// Done signals to stop watching for changes.
func (w *Watcher) Done() {
	w.once.Do(func() {
		close(w.done)
	})
}

// Err returns the watcher error, if any.
func (w *Watcher) Err() (err error) {
	w.mx.Lock()
	err = w.err
	w.mx.Unlock()

	return
}

// Watch begins watching for sample data updates.
// Interval is the interval between checking if there is an update available from Argus Monitor.
// If there is an update available, a new Sample is sent to the watchers output channel.
func (a *Api) Watch(interval time.Duration) *Watcher {
	changes := make(chan *Sample)
	watcher := &Watcher{
		C:    changes,
		done: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		defer close(changes)

		var cycle uint32

		for {
			select {
			case <-watcher.done:
				return
			case <-ticker.C:
			}

			currentCycle := a.CycleCounter()
			if cycle > 0 && cycle == currentCycle {
				continue
			}

			sample, err := a.Read()
			if err != nil {
				watcher.mx.Lock()
				watcher.err = err
				watcher.mx.Unlock()
				return
			}

			cycle = sample.CycleCounter

			select {
			case <-watcher.done:
				return
			case changes <- sample:
			}
		}
	}()

	return watcher
}

type Cached struct {
	api    *Api
	sample *Sample
}

// Read is cached read of the Argus Data Monitor API.
// Internally, the sample is only updated if there is no existing sample
// or if Argus Monitor hasn't changed since the last Read.
func (c *Cached) Read() (*Sample, error) {
	if c.sample != nil && c.sample.CycleCounter == c.api.CycleCounter() {
		return c.sample, nil
	}

	sample, err := c.api.Read()
	if err != nil {
		return nil, err
	}

	c.sample = sample
	return sample, nil
}

func (a *Api) Cached() *Cached {
	return &Cached{
		api: a,
	}
}
