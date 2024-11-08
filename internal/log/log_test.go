package log_test

import (
	"bytes"
	"order_manager/internal/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	errBuff := new(bytes.Buffer)
	logger := log.New(log.Debug, buf, errBuff)

	logger.Debugf("debug")
	logger.Infof("info")
	logger.Warningf("warning")
	logger.Errorf("error")

	assert.Contains(t, buf.String(), "debug")
	assert.Contains(t, buf.String(), "info")
	assert.NotContains(t, buf.String(), "warning")
	assert.NotContains(t, buf.String(), "error")

	assert.NotContains(t, errBuff.String(), "debug")
	assert.NotContains(t, errBuff.String(), "info")
	assert.Contains(t, errBuff.String(), "warning")
	assert.Contains(t, errBuff.String(), "error")
}

func TestInfoLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	errBuff := new(bytes.Buffer)
	logger := log.New(log.Info, buf, errBuff)

	logger.Debugf("debug")
	logger.Infof("info")
	logger.Warningf("warning")
	logger.Errorf("error")

	assert.NotContains(t, buf.String(), "debug")
	assert.Contains(t, buf.String(), "info")
	assert.NotContains(t, buf.String(), "warning")
	assert.NotContains(t, buf.String(), "error")

	assert.NotContains(t, errBuff.String(), "debug")
	assert.NotContains(t, errBuff.String(), "info")
	assert.Contains(t, errBuff.String(), "warning")
	assert.Contains(t, errBuff.String(), "error")
}

func TestWarningLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	errBuff := new(bytes.Buffer)
	logger := log.New(log.Warning, buf, errBuff)

	logger.Debugf("debug")
	logger.Infof("info")
	logger.Warningf("warning")
	logger.Errorf("error")

	assert.NotContains(t, buf.String(), "debug")
	assert.NotContains(t, buf.String(), "info")
	assert.NotContains(t, buf.String(), "warning")
	assert.NotContains(t, buf.String(), "error")

	assert.NotContains(t, errBuff.String(), "debug")
	assert.NotContains(t, errBuff.String(), "info")
	assert.Contains(t, errBuff.String(), "warning")
	assert.Contains(t, errBuff.String(), "error")
}

func TestErrorLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	errBuff := new(bytes.Buffer)
	logger := log.New(log.Error, buf, errBuff)

	logger.Debugf("debug")
	logger.Infof("info")
	logger.Warningf("warning")
	logger.Errorf("error")

	assert.NotContains(t, buf.String(), "debug")
	assert.NotContains(t, buf.String(), "info")
	assert.NotContains(t, buf.String(), "warning")
	assert.NotContains(t, buf.String(), "error")

	assert.NotContains(t, errBuff.String(), "debug")
	assert.NotContains(t, errBuff.String(), "info")
	assert.NotContains(t, errBuff.String(), "warning")
	assert.Contains(t, errBuff.String(), "error")
}
