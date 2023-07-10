package cloudtask

import (
	"context"
	"net/http"

	"google.golang.org/api/cloudtasks/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"metrio.net/fougere-lite/internal/common"
	"metrio.net/fougere-lite/internal/utils"
)

type Client struct {
	tasksService *cloudtasks.Service
}

func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	tasksService, err := cloudtasks.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		tasksService: tasksService,
	}, nil
}

func (c *Client) Create(config *Config) error {
	createChannel := make(chan common.Response, len(config.CloudTasks))
	for _, cloudtasks := range config.CloudTasks {
		go func(resp chan common.Response, cloudtasks TaskQueue) {
			_, err := c.get(cloudtasks.Name)
			if err != nil {
				if e, ok := err.(*googleapi.Error); ok && e.Code == http.StatusNotFound {
					utils.Logger.Debug("[%s] cloudtasks not found", cloudtasks.Name)

					if err := c.create(cloudtasks); err != nil {
						resp <- common.Response{Err: err}
						return
					}
				} else {
					utils.Logger.Errorf("[%s] error getting cloudtasks: %s", cloudtasks.Name, err)
					resp <- common.Response{Err: err}
					return
				}
			} else {
				if err := c.update(cloudtasks); err != nil {
					resp <- common.Response{Err: err}
					return
				}
			}
			resp <- common.Response{}
		}(createChannel, cloudtasks)
	}
	for range config.CloudTasks {
		resp := <-createChannel
		if resp.Err != nil {
			return resp.Err
		}
	}
	return nil
}

func (c *Client) get(name string) (*cloudtasks.Queue, error) {
	utils.Logger.Debug("[%s] getting cloudtasks", name)
	cloudtasks, err := c.tasksService.Projects.Locations.Queues.Get(name).Do()
	if err != nil {
		return nil, err
	}
	return cloudtasks, nil
}

func (c *Client) create(cloudtasks TaskQueue) error {
	utils.Logger.Infof("[%s] creating cloudtasks", cloudtasks.Name)
	spec := c.createQueueSpec(cloudtasks)
	_, err := c.tasksService.Projects.Locations.Queues.Create(cloudtasks.ProjectId, spec).Do()
	if err != nil {
		utils.Logger.Errorf("[%s] error creating cloudtasks: %s", spec.Name, err)
		return err
	}

	return nil
}

func (c *Client) update(cloudtasks TaskQueue) error {
	spec := c.createQueueSpec(cloudtasks)
	utils.Logger.Infof("[%s] updating cloudtasks", spec.Name)
	_, err := c.tasksService.Projects.Locations.Queues.Patch(spec.Name, spec).Do()
	if err != nil {
		utils.Logger.Errorf("[%s] error updating cloudtasks: %s", spec.Name, err)
		return err
	}
	return nil
}

func (c *Client) createQueueSpec(taskQueue TaskQueue) *cloudtasks.Queue {
	return &cloudtasks.Queue{
		Name: taskQueue.Name,
		RateLimits: &cloudtasks.RateLimits{
			MaxConcurrentDispatches: taskQueue.MaxConcurrentDispatches,
			MaxDispatchesPerSecond:  taskQueue.MaxDispatchesPerSecond,
		},
		RetryConfig: &cloudtasks.RetryConfig{
			MaxBackoff: taskQueue.MaxBackoff,
			MinBackoff: taskQueue.MinBackoff,
		},
	}
}
