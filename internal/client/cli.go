package client

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"metrio.net/fougere-lite/internal/gcp/cloudstorage"
	"metrio.net/fougere-lite/internal/gcp/cloudtask"
	"metrio.net/fougere-lite/internal/utils"
)

type ClientsCommand struct {
	cloudStorageClient *cloudstorage.Client
	cloudTasksClient   *cloudtask.Client
	clientConfigs      []ProductConfig
}

type ProductConfig struct {
	Client        string               `json:"client" validate:"omitempty"`
	StorageBucket *cloudstorage.Config `json:"storageBucket" validate:"omitempty,dive"`
	TasksQueue    *cloudtask.Config    `json:"tasksQueue" validate:"omitempty,dive"`
}

func NewClientsCommand() *cobra.Command {
	c := &ClientsCommand{}
	cmd := &cobra.Command{
		Use:   "clients",
		Short: "interacts with the GCP infrastructure components",
	}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create the client components",
		Run: func(cmd *cobra.Command, args []string) {
			c.getConfig()
			c.initClients()
			c.createClients()
		},
	}
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update the client components",
		Run: func(cmd *cobra.Command, args []string) {
			c.getConfig()
			c.initClients()
			c.updateClients()
		},
	}

	cmd.AddCommand(createCmd)
	cmd.AddCommand(updateCmd)
	return cmd
}

func (c *ClientsCommand) createClients() {
	for _, clientConfig := range c.clientConfigs {
		clientConfig := clientConfig
		if clientConfig.StorageBucket != nil {
			utils.CheckErr(cloudstorage.ValidateConfig(clientConfig.StorageBucket))
			if err := c.cloudStorageClient.Create(clientConfig.StorageBucket); err != nil {
				utils.CheckErr(err)
			}
		}
		if clientConfig.TasksQueue != nil {
			utils.CheckErr(cloudtask.ValidateConfig(clientConfig.TasksQueue))
			if err := c.cloudTasksClient.Create(clientConfig.TasksQueue); err != nil {
				utils.CheckErr(err)
			}
		}
	}
}

func (c *ClientsCommand) updateClients() {
	for _, clientConfig := range c.clientConfigs {
		clientConfig := clientConfig
		if clientConfig.StorageBucket != nil {
			utils.CheckErr(cloudstorage.ValidateConfig(clientConfig.StorageBucket))
			if err := c.cloudStorageClient.Create(clientConfig.StorageBucket); err != nil {
				utils.CheckErr(err)
			}
		}
		if clientConfig.TasksQueue != nil {
			utils.CheckErr(cloudtask.ValidateConfig(clientConfig.TasksQueue))
			if err := c.cloudTasksClient.Create(clientConfig.TasksQueue); err != nil {
				utils.CheckErr(err)
			}
		}
	}
}

func (c *ClientsCommand) initClients() error {
	ctx := context.Background()
	options := []option.ClientOption{}

	cloudStorageClient, err := cloudstorage.NewClient(ctx, options...)
	if err != nil {
		return err
	}
	c.cloudStorageClient = cloudStorageClient

	taskClient, err := cloudtask.NewClient(ctx, options...)
	if err != nil {
		return err
	}
	c.cloudTasksClient = taskClient

	return nil

}

func (c *ClientsCommand) getConfig() error {
	if !viper.InConfig("clients") {
		return fmt.Errorf("no clients config defined, please reference a fougere-lite.yaml")
	}
	clients := viper.GetStringMap("clients")
	for client := range clients {
		clientViper := viper.Sub(fmt.Sprintf("clients.%s", client))
		config := ProductConfig{
			Client: client,
		}
		storageConfig, err := cloudstorage.GetStorageConfig(clientViper, client)
		utils.CheckErr(err)
		config.StorageBucket = storageConfig
		taskConfig, err := cloudtask.GetTaskConfig(clientViper, client)
		utils.CheckErr(err)
		config.TasksQueue = taskConfig
		c.clientConfigs = append(c.clientConfigs, config)
	}
	return nil
}
