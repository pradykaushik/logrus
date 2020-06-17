package logrus

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

type testIOWriter struct {
	out io.Writer
}

func (w testIOWriter) Write(p []byte) (n int, err error) {
	type filtered struct {
		RequiredField1 string `json:"required_field1"`
		RequiredField2 string `json:"required_field2"`
		// Default fields.
		Msg string `json:"msg"`
		Level string `json:"level"`
		Timestamp string `json:"time"`
	}

	f := &filtered{}
	n = len(p)
	err = json.Unmarshal(p, f)
	if err != nil {
		return
	}
	p, err = json.Marshal(f)
	if err != nil {
		return
	}
	_, err = w.out.Write(p)
	return
}

func TestFilterFieldsToFiles(t *testing.T) {
	// creating a custom ioWriter.
	logFile, err := os.OpenFile("filtered_fields_file.log", os.O_CREATE|os.O_WRONLY, 0777)
	defer logFile.Close()
	assert.NoError(t, err)
	assert.NotNil(t, logFile)

	l := &Logger{
		Out:          io.MultiWriter(logFile, &testIOWriter{out: os.Stdout}),
		Formatter:    &JSONFormatter{},
		Level:        InfoLevel,
	}

	l.WithFields(Fields{
		"required_field1": "value1",
		"required_field2": "value2",
		"not_required_field3": "value3",
	}).Logln(InfoLevel, "Hello")
}