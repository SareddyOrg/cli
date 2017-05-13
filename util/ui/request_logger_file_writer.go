package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RequestLoggerFileWriter struct {
	ui        *UI
	lock      *sync.Mutex
	filePaths []string
	logFiles  []*os.File
}

func newRequestLoggerFileWriter(ui *UI, lock *sync.Mutex, filePaths []string) *RequestLoggerFileWriter {
	return &RequestLoggerFileWriter{
		ui:        ui,
		lock:      lock,
		filePaths: filePaths,
		logFiles:  []*os.File{},
	}
}

func (display *RequestLoggerFileWriter) DisplayBody(_ []byte) error {
	for _, logFile := range display.logFiles {
		_, err := logFile.WriteString(RedactedValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) DisplayDump(dump string) error {
	for _, logFile := range display.logFiles {
		_, err := logFile.WriteString(dump)
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) DisplayHeader(name string, value string) error {
	for _, logFile := range display.logFiles {
		_, err := logFile.WriteString(fmt.Sprintf("%s: %s\n", name, value))
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) DisplayHost(name string) error {
	for _, logFile := range display.logFiles {
		_, err := logFile.WriteString(fmt.Sprintf("Host: %s\n", name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) DisplayJSONBody(body []byte) error {
	if body == nil || len(body) == 0 {
		return nil
	}

	sanitized, err := SanitizeJSON(body)
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)
	encoder := json.NewEncoder(buff)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(sanitized)
	if err != nil {
		return err
	}

	for _, logFile := range display.logFiles {
		_, err = logFile.Write(buff.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) DisplayRequestHeader(method string, uri string, httpProtocol string) error {
	for _, logFile := range display.logFiles {
		_, err := logFile.WriteString(fmt.Sprintf("%s %s %s\n", method, uri, httpProtocol))
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) DisplayResponseHeader(httpProtocol string, status string) error {
	for _, logFile := range display.logFiles {
		_, err := logFile.WriteString(fmt.Sprintf("%s %s\n", httpProtocol, status))
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) DisplayType(name string, requestDate time.Time) error {
	for _, logFile := range display.logFiles {
		_, err := logFile.WriteString(fmt.Sprintf("%s: [%s]\n", name, requestDate.Format(time.RFC3339)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (display *RequestLoggerFileWriter) HandleInternalError(err error) {
	display.ui.DisplayWarning(err.Error())
}

func (display *RequestLoggerFileWriter) Start() error {
	display.lock.Lock()
	for _, filePath := range display.filePaths {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}

		logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return err
		}

		display.logFiles = append(display.logFiles, logFile)
	}
	return nil
}

func (display *RequestLoggerFileWriter) Stop() error {
	var err error

	for _, logFile := range display.logFiles {
		_, lastLineErr := logFile.WriteString("\n")
		closeErr := logFile.Close()
		switch {
		case closeErr != nil:
			err = closeErr
		case lastLineErr != nil:
			err = lastLineErr
		}
	}
	display.logFiles = []*os.File{}
	display.lock.Unlock()
	return err
}
