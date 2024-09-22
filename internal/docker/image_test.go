package docker

import (
	"fmt"
	"testing"

	"github.com/ChrisWiegman/kana-wordpress/internal/console"
	"github.com/ChrisWiegman/kana-wordpress/internal/docker/mocks"

	"github.com/docker/docker/api/types/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemoveImage(t *testing.T) {
	consoleOutput := new(console.Console)

	d, err := New(consoleOutput, "")
	assert.NoError(t, err)

	var tests = []struct {
		name                string
		imageDeleteResponse []image.DeleteResponse
		imageRemoveError    error
		expectedError       error
		expectedRemove      bool
	}{
		{
			"image doesn't exist to remove",
			[]image.DeleteResponse{},
			nil,
			nil,
			false},
		{
			"image successfully removed",
			[]image.DeleteResponse{
				{},
			},
			nil,
			nil,
			true},
		{
			"image successfully removed",
			[]image.DeleteResponse{
				{},
			},
			fmt.Errorf("image remove function hit error"),
			fmt.Errorf("image remove function hit error"),
			false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			apiClient := new(mocks.APIClient)
			apiClient.On("ImageRemove", mock.Anything, mock.Anything, mock.Anything).Return(test.imageDeleteResponse, test.imageRemoveError)

			d.apiClient = apiClient

			removed, err := d.removeImage("alpine")
			assert.Equal(t, test.expectedError, err, test.name)
			assert.Equal(t, test.expectedRemove, removed, test.name)
		})
	}
}
