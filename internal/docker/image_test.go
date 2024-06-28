package docker

import (
	"fmt"
	"testing"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker/mocks"

	"github.com/docker/docker/api/types/image"
	"github.com/knadh/koanf/v2"
	"github.com/moby/moby/pkg/jsonmessage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnsureImage(t *testing.T) {
	consoleOutput := new(console.Console)
	consoleOutput.JSON = true

	d, err := New(consoleOutput, "")
	assert.NoError(t, err)

	apiClient := new(mocks.APIClient)

	readCloser := &mocks.ReadCloser{
		ExpectedData: []byte(`{}`),
		ExpectedErr:  nil,
	}

	imageList := []image.Summary{
		{RepoTags: []string{
			"alpine:latest",
		}},
	}

	apiClient.On("ImagePull", mock.Anything, mock.Anything, mock.Anything).Return(readCloser, nil)
	apiClient.On("ImageList", mock.Anything, mock.Anything).Return(imageList, nil)

	d.apiClient = apiClient

	d.imageUpdateData = koanf.New(".")

	displayJSONMessagesStream = mocks.MockDisplayJSONMessagesStream
	mocks.MockedDisplayJSONMessagesStreamReturn = nil //nolint:gocritic

	err = d.EnsureImage("alpine", "", 1, consoleOutput)
	assert.Equal(t, nil, err)

	displayJSONMessagesStream = jsonmessage.DisplayJSONMessagesStream
}

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
		apiClient := new(mocks.APIClient)
		apiClient.On("ImageRemove", mock.Anything, mock.Anything, mock.Anything).Return(test.imageDeleteResponse, test.imageRemoveError)

		d.apiClient = apiClient

		removed, err := d.removeImage("alpine")
		assert.Equal(t, test.expectedError, err, test.name)
		assert.Equal(t, test.expectedRemove, removed, test.name)
	}
}
