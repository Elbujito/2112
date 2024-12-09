package handlers

import (
	"fmt"
	"strconv"

	"github.com/Elbujito/2112/src/templates/go-server/internal/config"
)

type TaskName string

type Task struct {
	Name         TaskName
	Description  string
	RequiredArgs []string
}

type TaskEnv struct {
	Env *config.SEnv
}

func (te *TaskEnv) SetEnv(env *config.SEnv) {
	te.Env = env
}

func ParseIntArg(args map[string]string, key string) (int, error) {
	value, ok := args[key]
	if !ok || value == "" {
		return 0, fmt.Errorf("missing required argument '%s'", key)
	}
	return strconv.Atoi(value)
}
