package handlers

import (
	"fmt"

	"github.com/Elbujito/2112/internal/config"
)

type Task struct {
	Name         string
	Description  string
	RequiredArgs []string
	Run          func(env *TaskEnv, args map[string]string) error
}

func (t *Task) Execute(args map[string]string) error {
	env := &TaskEnv{}
	env.SetEnv(config.Env)
	if err := validateArgs(args, t.RequiredArgs); err != nil {
		return err
	}
	if err := t.Run(env, args); err != nil {
		return err
	}
	return nil
}

func validateArgs(args map[string]string, requiredArgs []string) error {
	for _, arg := range requiredArgs {
		if val, ok := args[arg]; !ok || val == "" {
			return fmt.Errorf("missing required argument: %s", arg)
		}
	}
	return nil
}

type TaskEnv struct {
	Env *config.SEnv
}

func (te *TaskEnv) SetEnv(env *config.SEnv) {
	te.Env = env
}
