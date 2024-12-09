package celestrack

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Elbujito/2112/src/template/go-server/internal/api/mappers"
	"github.com/Elbujito/2112/src/template/go-server/internal/config"
)

type TestClient struct {
	env *config.SEnv
}

func NewTestClient(env *config.SEnv) *TestClient {
	return &TestClient{
		env: env,
	}
}

func (client *TestClient) FetchTest(ctx context.Context) ([]mappers.RawTest, error) {
	return []mappers.RawTest{}, fmt.Errorf("not implemented")
}

// parseFloat is a helper function to parse a float64 from a string.
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
