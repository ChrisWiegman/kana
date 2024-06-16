package site

import (
	"fmt"
	"testing"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/site/mocks"
	"github.com/stretchr/testify/assert"
)

func TestIsUsingSQLite(t *testing.T) {
	type databaseTest struct {
		Error          error
		Result         bool
		DockerResponse docker.ExecResult
		Description    string
	}

	tests := []databaseTest{
		{
			Description: "Test without error but response is not `true`",
			Error:       nil,
			Result:      false,
			DockerResponse: docker.ExecResult{
				StdOut:   "sqlite",
				StdErr:   "",
				ExitCode: 0},
		},
		{
			Description: "Test response is `true` but we get an error",
			Error:       fmt.Errorf("error"),
			Result:      false,
			DockerResponse: docker.ExecResult{
				StdOut:   "sqlite",
				StdErr:   "this is an error",
				ExitCode: 1},
		},
		{
			Description: "Test response is `true` without error",
			Error:       nil,
			Result:      true,
			DockerResponse: docker.ExecResult{
				StdOut:   "true",
				StdErr:   "",
				ExitCode: 0},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			mockCli := mocks.NewCli(t)

			mockCli.On("WordPress", "echo $KANA_SQLITE", false, false).Return(test.DockerResponse, test.Error)

			s := Site{
				Cli: mockCli,
			}

			result, err := s.isUsingSQLite()
			assert.Equal(t, test.Result, result)
			assert.Equal(t, test.Error, err)
		})
	}
}

func TestVerifyDatabase(t *testing.T) {
	type databaseTest struct {
		Error       error
		Code        int64
		Description string
		Return      error
	}

	tests := []databaseTest{
		{
			Description: "Test without error",
			Error:       nil,
			Code:        0,
			Return:      nil,
		},
		{
			Description: "Test without error but exit code is not 0",
			Error:       nil,
			Code:        1,
			Return:      fmt.Errorf("database verification failed"),
		},
		{
			Description: "Response code is 0 but we get an error",
			Error:       fmt.Errorf("error"),
			Code:        0,
			Return:      fmt.Errorf("database verification failed"),
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			mockCli := mocks.NewCli(t)

			mockCli.On("WPCli", []string{"db", "check"}, false, &console.Console{}).Return(test.Code, "", test.Error)

			s := Site{
				Cli:                    mockCli,
				maxVerificationRetries: 3,
			}

			err := s.verifyDatabase(&console.Console{})
			assert.Equal(t, test.Return, err)
		})
	}
}
