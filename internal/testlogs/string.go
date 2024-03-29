package testlogs

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

// ToString transforms a slice of logrus entry into a string concatenation.
func ToString(entries []*logrus.Entry) string {
	var agg bytes.Buffer
	for _, entry := range entries {
		b, _ := entry.Bytes()
		agg.Write(b)
		agg.WriteString("\n")
	}
	return agg.String()
}
