package network

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionSuite))
}

type ConnectionSuite struct {
	suite.Suite
}
