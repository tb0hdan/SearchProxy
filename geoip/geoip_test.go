package geoip

import (
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := testifyAssert.New(t)
	assert.IsType(&DB{}, New(""))
}

func TestDB_DistanceIP(t *testing.T) {
	assert := testifyAssert.New(t)
	geo := New("../GeoLite2-City.mmdb")
	dist, err := geo.DistanceIP("1.2.3.4", "4.4.4.4")
	assert.Equal(8795, int(dist))
	assert.Nil(err)
}

func TestDB_DistanceIPLatLon(t *testing.T) {
	assert := testifyAssert.New(t)
	geo := New("../GeoLite2-City.mmdb")
	dist, err := geo.DistanceIPLatLon("1.2.3.4", 50.5, 30.5)
	assert.Equal(751, int(dist))
	assert.Nil(err)
}

func TestDB_DistanceLatLon(t *testing.T) {
	assert := testifyAssert.New(t)
	geo := New("../GeoLite2-City.mmdb")
	dist := geo.DistanceLatLon(120.1, -30.2, 50.5, 30.5)
	assert.Equal(6587, int(dist))
}

func TestDB_GeoIPLookupIP(t *testing.T) {
	assert := testifyAssert.New(t)
	geo := New("../GeoLite2-City.mmdb")
	info, err := geo.LookupIP("200.100.50.25")
	assert.Equal("São Paulo", info.CityName)
	assert.Nil(err)
}

func TestDB_LookupDomain(t *testing.T) {
	assert := testifyAssert.New(t)
	geo := New("../GeoLite2-City.mmdb")
	info, err := geo.LookupDomain("example.com")
	assert.Equal("Norwell", info.CityName)
	assert.Nil(err)
}

func TestDB_LookupIP(t *testing.T) {
	assert := testifyAssert.New(t)
	geo := New("../GeoLite2-City.mmdb")
	info, err := geo.LookupIP("125.200.125.200")
	assert.Equal("Minokamo", info.CityName)
	assert.Nil(err)
}

func TestDB_LookupURL(t *testing.T) {
	assert := testifyAssert.New(t)
	geo := New("../GeoLite2-City.mmdb")
	info, err := geo.LookupURL("https://www.foobar.com/1/2/3")
	assert.Equal("El Segundo", info.CityName)
	assert.Nil(err)
}
