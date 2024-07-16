package validator

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (*Validator, repository.MockDatabase) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockDB := repository.NewMockDatabase(ctrl)
	mockValidatorCommand := NewValidator(mockDB)

	return mockValidatorCommand, *mockDB
}

func TestImport(t *testing.T) {
	validatorCmd, mockDB := setup(t)

	t.Run("normal", func(t *testing.T) {
		defer gock.Off()
		gock.New("http://foo.com").
			Get("/bar").
			Reply(200).
			BodyString("Name,Email\nValidator1,validator1@abc.com\nValidator2,validator2@abc.com\nValidator3,validator3@abc.com")

		cmd := &command.Command{}

		mockDB.EXPECT().AddValidator(gomock.Any()).Return(nil).AnyTimes()

		args := make(map[string]string)
		args["file"] = "http://foo.com/bar"
		result := validatorCmd.importHandler(nil, cmd, args)

		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Validators created successfully!")
	})
}
