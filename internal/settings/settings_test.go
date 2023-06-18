package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettings(t *testing.T) {
	settings, err := NewSettings()

	assert.False(t, settings.Activate)
	assert.Equal(t, nil, err)
}
