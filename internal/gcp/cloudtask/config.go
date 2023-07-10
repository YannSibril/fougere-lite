// ©Copyright 2022-2023 Metrio
package cloudtask

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	CloudTasks map[string]TaskQueue `mapstructure:"cloudTasks" validate:"dive"`
}

// TaskQueue contains the information required to create a Cloud Tasks Queue in gcp.
// A cloud task queue is used to store all kinds of task.
type TaskQueue struct {
	Name                    string  `json:"name" validate:"required"`
	Region                  string  `json:"region" validate:"required"`
	ProjectId               string  `json:"projectId" validate:"required"`
	MinBackoff              string  `json:"minBackoff" validate:"required"`
	MaxBackoff              string  `json:"maxBackoff" validate:"required"`
	MaxConcurrentDispatches int64   `json:"MaxConcurrentDispatches" validate:"required"`
	MaxDispatchesPerSecond  float64 `json:"MaxDispatchesPerSecond" validate:"required"`
	ClientName              string
}

func GetTaskConfig(viperConfig *viper.Viper, clientName string) (*Config, error) {
	if viperConfig == nil {
		return nil, nil
	}

	var taskConfig Config
	err := viperConfig.Unmarshal(&taskConfig)
	if err != nil {
		return nil, err
	}

	for name, queue := range taskConfig.CloudTasks {
		queue.Name = strings.Join([]string{clientName, name, queue.ProjectId}, "-")
		if err != nil {
			return nil, err
		}
		queue.ClientName = clientName

		taskConfig.CloudTasks[name] = queue
	}
	return &taskConfig, nil
}

func ValidateConfig(config *Config) error {
	v := validator.New()
	if err := v.Struct(config); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s validate failed on the %s rule", err.Namespace(), err.Tag())
		}
	}
	return nil
}
