package pipeline

type FileWriter struct {
}

func NewFileWriter() *FileWriter {
	w := &FileWriter{}
	return w
}

func (w *FileWriter) Write() {
}
