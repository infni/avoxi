package test

import (
	"ipauthorize/internal/pkg/countycodes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTestComparer() countycodes.CountryCodeComparer {
	comparer := countycodes.NewContryCodeComparer("../GeoLite2-Country.mmdb")

	return comparer
}

func TestComparerNilIp(t *testing.T) {

	comparer := NewTestComparer()

	actual, acturalErr := comparer.IsInCountry("", nil)

	assert.Error(t, acturalErr)
	assert.False(t, actual)
}

func TestComparerBadIp(t *testing.T) {

	comparer := NewTestComparer()

	actual, acturalErr := comparer.IsInCountry("10.", nil)

	assert.Error(t, acturalErr)
	assert.False(t, actual)
}

func TestComparerNilCountries(t *testing.T) {

	comparer := NewTestComparer()

	actual, acturalErr := comparer.IsInCountry("152.216.7.110", nil)

	assert.Nil(t, acturalErr)
	assert.False(t, actual)
}

func TestComparerEmptyCountries(t *testing.T) {

	comparer := NewTestComparer()

	actual, acturalErr := comparer.IsInCountry("152.216.7.110", []string{})

	assert.Nil(t, acturalErr)
	assert.False(t, actual)
}

func TestComparerNoMatchCounties(t *testing.T) {

	comparer := NewTestComparer()

	actual, acturalErr := comparer.IsInCountry("152.216.7.110", []string{"one", "two"})

	assert.Nil(t, acturalErr)
	assert.False(t, actual)
}

func TestComparerMatchingCountry(t *testing.T) {

	comparer := NewTestComparer()

	actual, acturalErr := comparer.IsInCountry("152.216.7.110", []string{"United States", "Uganda"})

	assert.Nil(t, acturalErr)
	assert.True(t, actual)
}
