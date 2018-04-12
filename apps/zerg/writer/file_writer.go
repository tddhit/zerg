package writer

import (
	"bytes"
	"os"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/types"
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

func (w *FileWriter) Write(item *types.Item) {
	var buf bytes.Buffer
	count := 0
	for key, value := range item.Dict {
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(value.(string))
		buf.WriteString("\t")
		count++
	}
	if count > 0 {
		buf.WriteString("\n")
	}
	w.WriteString(buf.String())
}
