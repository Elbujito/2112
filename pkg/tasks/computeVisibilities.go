package tasks

func init() {
	task := &Task{
		Name:         "computeVisibilities",
		Description:  "Computes the visibilities for all tiles given tles and satellites",
		RequiredArgs: []string{},
		Run:          execComputeVisibilitiesTask,
	}
	Tasks.AddTask(task)
}

func execComputeVisibilitiesTask(env *TaskEnv, args map[string]string) error {

	return nil
}
