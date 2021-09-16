// Telegraf execd compatible plugin for Argus Monitor
package main

import (
	"bufio"
	"fmt"
	"github.com/relvacode/argus"
	"os"
	"strings"
	"time"
)

func safe(v string) string {
	var s strings.Builder
	for _, ch := range v {
		switch ch {
		case ',', ' ', '=':
			s.WriteRune('\\')
			s.WriteRune(ch)
		default:
			s.WriteRune(ch)
		}
	}

	return s.String()
}

func write(m argus.Measurement, timestamp int64) error {
	_, err := fmt.Fprintf(os.Stdout, "%s,label=%s,sensor=%d value=%f %d\n", safe(m.Type.String()), safe(m.Label), m.SensorIndex, m.Value, timestamp)
	return err
}

func Main() error {
	api, err := argus.Open()
	if err != nil {
		return err
	}

	defer api.Close()

	cached := api.Cached()

	notify := bufio.NewScanner(os.Stdin)
	for notify.Scan() {
		sample, err := cached.Read()
		if err != nil {
			return err
		}

		timestamp := time.Now().UTC().UnixNano()

		for _, m := range sample.Data {
			err = write(m, timestamp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	err := Main()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
