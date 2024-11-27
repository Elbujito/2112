package tasks

import (
	"errors"
	"testing"
	"time"

	"github.com/Elbujito/2112/internal/api/mappers"
	"github.com/Elbujito/2112/internal/data/models"
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

func (m *MockTLEService) Upsert(tle *models.TLE) error {
	args := m.Called(tle)
	return args.Error(0)
}

func (m *MockTLEService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestTLEHandler_Run(t *testing.T) {
	tests := []struct {
		name            string
		mockFetchTLE    func(string) ([]*mappers.RawTLE, error)
		mockSatellite   func(*MockSatelliteService)
		mockTLE         func(*MockTLEService)
		category        string
		expectedErr     bool
		expectedUpdates int
	}{
		{
			name: "happy path - updates existing TLE",
			mockFetchTLE: func(category string) ([]*mappers.RawTLE, error) {
				return []*mappers.RawTLE{
					{
						NoradID: "25544",
						Line1:   "1 25544U 98067A   21273.00000000  .00001264  00000-0  29647-4 0  9998",
						Line2:   "2 25544  51.6441 245.2066 0003157  97.5202 262.6127 15.48907224281129",
					},
				}, nil
			},
			mockSatellite: func(s *MockSatelliteService) {
				s.On("FindByNoradID", "25544").Return(&models.Satellite{NoradID: "25544", Name: "ISS"}, nil)
			},
			mockTLE: func(t *MockTLEService) {
				t.On("FindByNoradID", "25544").Return([]*models.TLE{
					{
						NoradID: "25544",
						Line1:   "old-line1",
						Line2:   "old-line2",
						Epoch:   time.Now(),
					},
				}, nil)
				t.On("Upsert", mock.AnythingOfType("*models.TLE")).Return(nil)
			},
			category:        "active",
			expectedErr:     false,
			expectedUpdates: 1,
		},
		{
			name: "error in fetchTLEHandler",
			mockFetchTLE: func(category string) ([]*mappers.RawTLE, error) {
				return nil, errors.New("fetch error")
			},
			mockSatellite:   func(s *MockSatelliteService) {},
			mockTLE:         func(t *MockTLEService) {},
			category:        "active",
			expectedErr:     true,
			expectedUpdates: 0,
		},
		{
			name: "satellite not found - creates new satellite",
			mockFetchTLE: func(category string) ([]*mappers.RawTLE, error) {
				return []*mappers.RawTLE{
					{
						NoradID: "12345",
						Line1:   "1 12345U 98067A   21273.00000000  .00001264  00000-0  29647-4 0  9998",
						Line2:   "2 12345  51.6441 245.2066 0003157  97.5202 262.6127 15.48907224281129",
					},
				}, nil
			},
			mockSatellite: func(s *MockSatelliteService) {
				s.On("FindByNoradID", "12345").Return(nil, nil)
				s.On("Save", mock.AnythingOfType("*models.Satellite")).Return(nil)
			},
			mockTLE: func(t *MockTLEService) {
				t.On("FindByNoradID", "12345").Return([]*models.TLE{}, nil)
				t.On("Upsert", mock.AnythingOfType("*models.TLE")).Return(nil)
			},
			category:        "active",
			expectedErr:     false,
			expectedUpdates: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSatellite := new(MockSatelliteService)
			mockTLE := new(MockTLEService)
			tt.mockSatellite(mockSatellite)
			tt.mockTLE(mockTLE)

			handler := NewTLEHandler(mockSatellite, mockTLE, tt.mockFetchTLE)

			err := handler.Run(tt.category)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockSatellite.AssertExpectations(t)
			mockTLE.AssertExpectations(t)
		})
	}
}
