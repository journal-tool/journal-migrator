package lib

import (
	"github.com/stretchr/testify/suite"
	"journal-migrator/lib"
	"testing"
)

func TestSliceFunctions(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping unit tests")
	}

	suite.Run(t, new(SliceFunctionsSuite))
}

type SliceFunctionsSuite struct {
	suite.Suite
}

func (suite *SliceFunctionsSuite) TestIntersect() {
	var argument [][]string
	var expected []string

	argument = [][]string{{"a", "b"}}
	expected = []string{"a", "b"}
	suite.Equal(expected, lib.Intersect(argument...))

	argument = [][]string{{"a", "b"}, {"b", "c"}}
	expected = []string{"b"}
	suite.Equal(expected, lib.Intersect(argument...))

	argument = [][]string{{"a", "b"}, {"c", "d"}}
	expected = nil
	suite.Equal(expected, lib.Intersect(argument...))

	argument = [][]string{{"a", "b"}, {}}
	expected = nil
	suite.Equal(expected, lib.Intersect(argument...))

	argument = [][]string{}
	expected = nil
	suite.Equal(expected, lib.Intersect(argument...))
}
