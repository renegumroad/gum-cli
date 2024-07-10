package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type logSuite struct {
	suite.Suite
}

func (s *logSuite) TestErrorLevel() {
	err := Initialize("invalid")

	assert.Error(s.T(), err)
}

func TestLogSuite(t *testing.T) {
	suite.Run(t, new(logSuite))
}
