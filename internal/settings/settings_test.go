package settings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettings(t *testing.T) {
	settings, err := NewSettings()

	fmt.Println(settings.Activate)

	assert.False(t, settings.Activate)
	assert.Equal(t, nil, err)
}
