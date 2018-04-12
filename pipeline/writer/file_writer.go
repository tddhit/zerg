package writer

import (
	"bytes"
	"os"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/types"
)

type FileWriter struct {
	*os.File
}

func NewFileWriter(filePath string) *FileWriter {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Errorf("failed open file: %s, %s", filePath, err)
		return nil
	}
	w := &FileWriter{
		File: file,
	}
	return w
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
