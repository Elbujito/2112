package handlers

import (
	"github.com/Elbujito/2112/internal/config"
)

type Task struct {
	Name         string
	Description  string
	RequiredArgs []string
}

type TaskEnv struct {
	Env *config.SEnv
}

func (te *TaskEnv) SetEnv(env *config.SEnv) {
	te.Env = env
}
