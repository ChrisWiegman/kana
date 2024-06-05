package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettings(t *testing.T) {
	settings, err := NewSettings("1.0.0")

	assert.NotNil(t, settings)
	assert.False(t, settings.Activate)
	assert.Equal(t, nil, err)
}
