package homebrew

import (
	"github.com/renegumroad/gum-cli/internal/homebrew/mockhomebrew"
	"github.com/stretchr/testify/suite"
)

type homebrewSuite struct {
	suite.Suite
	mockBrew *mockhomebrew.MockClient
}
