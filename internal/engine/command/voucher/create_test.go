package voucher

import (
	"github.com/pagu-project/Pagu/internal/repository"
	"go.uber.org/mock/gomock"
)

func setup() { //nolint
	_ = repository.NewMockDatabase(gomock.NewController(nil))
}
