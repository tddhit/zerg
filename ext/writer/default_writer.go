package writer

import (
	"bytes"
	"os"
	"sort"

	"github.com/tddhit/zerg"
)

type DefaultWriter struct {
	name string
}

func NewDefaultWriter() *DefaultWriter {
	return &DefaultWriter{}
}

func (w *DefaultWriter) Name() string {
	return "DEFAULT_WRITER"
}

func (w *DefaultWriter) Write(item *zerg.Item) {
	var (
		buf        bytes.Buffer
		sortedKeys = make([]string, 0, len(item.Dict))
	)
	for key, _ := range item.Dict {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	for _, key := range sortedKeys {
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(item.Dict[key].(string))
		buf.WriteString("\t")
	}
	if len(sortedKeys) > 0 {
		buf.WriteString("\n")
	}
	os.Stdout.WriteString(buf.String())
}
