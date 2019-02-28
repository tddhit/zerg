package writer

import (
	"bytes"
	"os"
	"sort"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg"
)

type FileWriter struct {
	*os.File
	name string
}

func NewFileWriter(name, filePath string) *FileWriter {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error("failed open file: %s, %s", filePath, err)
		return nil
	}
	w := &FileWriter{
		File: file,
		name: name,
	}
	return w
}

func (w *FileWriter) Name() string {
	return w.name
}

func (w *FileWriter) Write(item *zerg.Item) {
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
	w.WriteString(buf.String())
}
