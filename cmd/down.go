/*
MIT License

Copyright (c) Nhost

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"context"
	"io"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mrinalwahal/cli/nhost"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove local Nhost backend started by \"nhost dev\"",
	Long:  "Stop and remove local Nhost backend started by \"nhost dev\".",
	Run: func(cmd *cobra.Command, args []string) {

		// connect to docker client
		ctx := context.Background()
		docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal("Failed to connect to docker client")
		}

		if err := shutdownServices(docker, ctx, LOG_FILE); err != nil {
			log.Error("Failed to shut down Nhost services")
		}

		deletePath(path.Join(nhost.DOT_NHOST, "Dockerfile-api"))

		if contains(args, "exit") {
			log.Info("Cleanup complete. See you later, grasshopper!")
			os.Exit(0)
		}
	},
}

func shutdownServices(client *client.Client, ctx context.Context, logFile string) error {

	// get running containers with prefix "nhost_"
	containers, err := getContainers(client, ctx, "nhost")
	if err != nil {
		return err
	}

	if len(containers) > 0 {
		log.Info("Running Nhost services detected. Shutting them down.")
	}

	for _, container := range containers {

		if logFile != "" {

			// generate container logs and write them to logFile
			if err = writeContainerLogs(client, ctx, logFile, container); err != nil {
				return err
			}
		}

		// stop the container
		if err = stopContainer(client, ctx, container); err != nil {
			return err
		}

		// remove all running containers with prefix "nhost_"
		if err = removeContainer(client, ctx, container); err != nil {
			return err
		}
	}

	// search and delete Nhost network too,
	// if you don't do this,
	// docker will get confused about ambigous network names

	network, err := getNetwork(client, ctx, "nhost")
	if err != nil {
		return err
	}

	if network != "" {
		err = removeNetwork(client, ctx, network)
		return err

	}
	return err
}

// returns the list of running containers whose names have specified prefix
func getContainers(cli *client.Client, ctx context.Context, prefix string) ([]types.Container, error) {

	log.Debug("Fetching running containers with names having the prefix: ", prefix)

	var response []types.Container
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	for _, container := range containers {
		if strings.Contains(container.Names[0], prefix) {
			response = append(response, container)
		}
	}

	return response, err
}

// removes a given network by ID
func removeNetwork(cli *client.Client, ctx context.Context, ID string) error {

	log.Debug("Removing network: ", ID)

	err := cli.NetworkRemove(ctx, ID)
	return err
}

// fetches ID of docker network by name
func getNetwork(cli *client.Client, ctx context.Context, name string) (string, error) {

	log.WithField("component", name).Debug("Fetching network")

	response, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", err
	}

	for _, network := range response {
		if network.Name == name {
			return network.ID, nil
		}
	}

	return "", err
}

// restarts given container
func restartContainer(cli *client.Client, ctx context.Context, container types.Container) error {

	log.WithField("component", container.Names[0]).Debug("Restarting container")

	err := cli.ContainerStop(ctx, container.ID, nil)
	return err
}

// stops given container
func stopContainer(cli *client.Client, ctx context.Context, container types.Container) error {

	log.WithField("component", container.Names[0]).Debug("Stopping container")

	err := cli.ContainerStop(ctx, container.ID, nil)
	return err
}

// fetches the logs of a specific container
// and writes them to a log file
func writeContainerLogs(cli *client.Client, ctx context.Context, filePath string, container types.Container) error {

	log.WithField("component", container.Names[0]).Debug("Writing container logs to ", filePath)

	options := types.ContainerLogsOptions{ShowStdout: true}

	out, err := cli.ContainerLogs(ctx, container.ID, options)
	if err != nil {
		return err
	}

	// write the fetched logs to a file
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// handle error
	defer f.Close()

	_, err = io.Copy(f, out)
	return err
}

// removes given container
func removeContainer(cli *client.Client, ctx context.Context, container types.Container) error {

	log.WithField("component", container.Names[0]).Debug("Removing container")

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	err := cli.ContainerRemove(ctx, container.ID, removeOptions)
	return err
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
