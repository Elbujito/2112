package handlers

import (
	"github.com/Elbujito/2112/internal/config"
)

type Task struct {
	Name         string
	Description  string
	RequiredArgs []string
}

// func validateArgs(args map[string]string, requiredArgs []string) error {
// 	for _, arg := range requiredArgs {
// 		if val, ok := args[arg]; !ok || val == "" {
// 			return fmt.Errorf("missing required argument: %s", arg)
// 		}
// 	}
// 	return nil
// }

type TaskEnv struct {
	Env *config.SEnv
}

func (te *TaskEnv) SetEnv(env *config.SEnv) {
	te.Env = env
}
