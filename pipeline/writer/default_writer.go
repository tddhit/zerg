package writer

import (
	"bytes"
	"os"

	"github.com/tddhit/zerg/types"
)

type ConsoleWriter struct {
	name string
}

func NewConsoleWriter(name string) *ConsoleWriter {
	w := &ConsoleWriter{
		name: name,
	}
	return w
}

func (w *ConsoleWriter) Name() string {
	return w.name
}

func (w *ConsoleWriter) Write(item *types.Item) {
	var buf bytes.Buffer
	count := 0
	for key, value := range item.Dict {
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(value)
		buf.WriteString("\t")
		count++
	}
	if count > 0 {
		buf.WriteString("\n")
	}
	os.Stdout.WriteString(buf.String())
}
