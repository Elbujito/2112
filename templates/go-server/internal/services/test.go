package services

import "github.com/Elbujito/2112/template/go-server/internal/domain"

type TestService struct {
	repo domain.TestRepository
}

// NewTestService creates a new instance of TestService.
func NewTestService(testRepo domain.TestRepository) TestService {
	return TestService{repo: testRepo}
}
