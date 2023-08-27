package docker

import (
	"fmt"
	"testing"
	"time"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker/mocks"
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/pkg/jsonmessage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnsureImage(t *testing.T) {
	consoleOutput := new(console.Console)
	consoleOutput.JSON = true

	d, err := NewDockerClient(consoleOutput, "")
	assert.NoError(t, err)

	apiClient := new(mocks.APIClient)

	readCloser := &mocks.ReadCloser{
		ExpectedData: []byte(`{}`),
		ExpectedErr:  nil,
	}

	imageList := []types.ImageSummary{
		{RepoTags: []string{
			"alpine:latest",
		}},
	}

	apiClient.On("ImagePull", mock.Anything, mock.Anything, mock.Anything).Return(readCloser, nil)
	apiClient.On("ImageList", mock.Anything, mock.Anything).Return(imageList, nil)

	d.apiClient = apiClient

	viper := new(mocks.ViperClient)
	viper.On("ReadInConfig").Return(nil)
	viper.On("GetTime", mock.Anything).Return(time.Now())
	viper.On("Set", mock.Anything, mock.Anything).Return()
	viper.On("WriteConfig").Return(nil)

	d.imageUpdateData = viper

	displayJSONMessagesStream = mocks.MockDisplayJSONMessagesStream
	mocks.MockedDisplayJSONMessagesStreamReturn = nil //nolint:gocritic

	err = d.EnsureImage("alpine", 1, consoleOutput)
	assert.Equal(t, nil, err)

	displayJSONMessagesStream = jsonmessage.DisplayJSONMessagesStream
}

func TestRemoveImage(t *testing.T) {
	consoleOutput := new(console.Console)

	d, err := NewDockerClient(consoleOutput, "")
	assert.NoError(t, err)

	var tests = []struct {
		name                string
		imageDeleteResponse []types.ImageDeleteResponseItem
		imageRemoveError    error
		expectedError       error
		expectedRemove      bool
	}{
		{
			"image doesn't exist to remove",
			[]types.ImageDeleteResponseItem{},
			nil,
			nil,
			false},
		{
			"image successfully removed",
			[]types.ImageDeleteResponseItem{
				{},
			},
			nil,
			nil,
			true},
		{
			"image successfully removed",
			[]types.ImageDeleteResponseItem{
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
