package tasks

import (
	"testing"
	"time"

	"github.com/Elbujito/2112/pkg/api/mappers"
	"github.com/Elbujito/2112/pkg/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of SatelliteService
type MockSatelliteService struct {
	mock.Mock
}

func (m *MockSatelliteService) FindByNoradID(noradID string) (*models.Satellite, error) {
	args := m.Called(noradID)
	return args.Get(0).(*models.Satellite), args.Error(1)
}

func (m *MockSatelliteService) Find(id string) (*models.Satellite, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Satellite), args.Error(1)
}

func (m *MockSatelliteService) FindAll() ([]*models.Satellite, error) {
	args := m.Called()
	return args.Get(0).([]*models.Satellite), args.Error(1)
}

func (m *MockSatelliteService) Save(satellite *models.Satellite) error {
	args := m.Called(satellite)
	return args.Error(0)
}

func (m *MockSatelliteService) Update(satellite *models.Satellite) error {
	args := m.Called(satellite)
	return args.Error(0)
}

func (m *MockSatelliteService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock implementation of TLEService
type MockTLEService struct {
	mock.Mock
}

func (m *MockTLEService) FindByNoradID(noradID string) ([]*models.TLE, error) {
	args := m.Called(noradID)
	return args.Get(0).([]*models.TLE), args.Error(1)
}

func (m *MockTLEService) FindAll() ([]*models.TLE, error) {
	args := m.Called()
	return args.Get(0).([]*models.TLE), args.Error(1)
}

func (m *MockTLEService) Save(tle *models.TLE) error {
	args := m.Called(tle)
	return args.Error(0)
}

func (m *MockTLEService) Update(tle *models.TLE) error {
	args := m.Called(tle)
	return args.Error(0)
}

func (m *MockTLEService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func TestFetchAndUpsertTLEs_HappyPath(t *testing.T) {
	// Arrange
	mockSatellite := new(MockSatelliteService)
	mockTLE := new(MockTLEService)

	mockFetchTLE := func(category string) ([]*mappers.RawTLE, error) {
		return []*mappers.RawTLE{
			{
				NoradID: "25544",
				Line1:   "1 25544U 98067A   21273.00000000  .00001264  00000-0  29647-4 0  9998",
				Line2:   "2 25544  51.6441 245.2066 0003157  97.5202 262.6127 15.48907224281129",
			},
		}, nil
	}

	foundSatellite := &models.Satellite{
		NoradID: "25544",
		Name:    "ISS (ZARYA)",
	}
	mockSatellite.On("FindByNoradID", "25544").Return(foundSatellite, nil)

	// Mock Save for Satellite
	mockSatellite.On("Save", mock.MatchedBy(func(sat *models.Satellite) bool {
		t.Logf("Save called with Satellite: %+v", sat) // Debug log
		return sat.NoradID == "25544" && sat.Name == "ISS (ZARYA)"
	})).Return(nil)

	existingTLE := &models.TLE{
		NoradID: "25544",
		Line1:   "1 25544U 98067A   21273.00000000  .00001264  00000-0  29647-4 0  9998",
		Line2:   "2 25544  51.6441 245.2066 0003157  97.5202 262.6127 15.48907224281129",
		Epoch:   time.Date(2021, time.September, 30, 0, 0, 0, 0, time.UTC),
	}
	mockTLE.On("FindByNoradID", "25544").Return([]*models.TLE{existingTLE}, nil)

	// Mock Save for TLE
	mockTLE.On("Save", mock.MatchedBy(func(tle *models.TLE) bool {
		t.Logf("Save called with TLE: %+v", tle) // Debug log
		return tle.NoradID == "25544" &&
			tle.Line1 == "1 25544U 98067A   21273.00000000  .00001264  00000-0  29647-4 0  9998" &&
			tle.Line2 == "2 25544  51.6441 245.2066 0003157  97.5202 262.6127 15.48907224281129"
	})).Return(nil)

	// Mock Update behavior
	mockTLE.On("Update", mock.MatchedBy(func(tle *models.TLE) bool {
		t.Logf("Update called with TLE: %+v", tle) // Debug log
		return tle.NoradID == "25544" &&
			tle.Line1 == "1 25544U 98067A   21273.00000000  .00001264  00000-0  29647-4 0  9998" &&
			tle.Line2 == "2 25544  51.6441 245.2066 0003157  97.5202 262.6127 15.48907224281129"
	})).Return(nil)

	// Act
	err := fetchAndUpsertTLEs("active", mockSatellite, mockTLE, mockFetchTLE)

	// Assert
	assert.NoError(t, err)
	mockSatellite.AssertCalled(t, "FindByNoradID", "25544")
	// mockSatellite.AssertCalled(t, "Save", mock.Anything)
	mockTLE.AssertCalled(t, "FindByNoradID", "25544")
	// mockTLE.AssertCalled(t, "Save", mock.Anything)
	mockTLE.AssertCalled(t, "Update", mock.Anything)
}
