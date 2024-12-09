package handlers

import (
	"context"

	"github.com/Elbujito/2112/src/template/go-server/internal/domain"
)

type TestHandler struct {
	testRepo domain.TestRepository
}

func NewTestHandler(
	satelliteRepo domain.TestRepository) TestHandler {
	return TestHandler{
		testRepo: satelliteRepo,
	}
}

func (h *TestHandler) GetTask() Task {
	return Task{
		Name:         "test_handler",
		Description:  "simple template for handler",
		RequiredArgs: []string{""},
	}
}

func (h *TestHandler) Run(ctx context.Context, args map[string]string) error {

	return nil
}
