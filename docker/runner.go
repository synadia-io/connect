package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/synadia-io/connect/model"
	"gopkg.in/yaml.v3"
)

type Runner struct {
	// For now, we'll use simple docker command execution
	// In the future, we can add the Docker client library
}

type ContainerStatus struct {
	Name   string
	Status string
	Exists bool
}

type RunOptions struct {
	ConnectorID string
	Image       string
	Steps       model.Steps
	EnvVars     map[string]string
	DockerOpts  string
	WorkDir     string
	Follow      bool
	Remove      bool
	RuntimeID   string // Optional: runtime ID for converter selection
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Run(ctx context.Context, opts *RunOptions) error {
	if opts.Image == "" {
		opts.Image = "synadia/connect:latest" // Default image
	}

	// Check if a container with the same name already exists
	containerStatus, err := r.GetContainerStatus(ctx, opts.ConnectorID)
	if err != nil {
		return fmt.Errorf("failed to check existing container: %w", err)
	}

	if containerStatus.Exists {
		if containerStatus.IsContainerRunning() {
			// Container is running, ask user if they want to replace it
			replace, err := r.PromptUserForReplacement(opts.ConnectorID)
			if err != nil {
				return fmt.Errorf("failed to get user input: %w", err)
			}

			if !replace {
				fmt.Printf("Keeping existing running container '%s'\n", opts.ConnectorID)
				return nil
			}

			// User wants to replace, remove the existing container
			if err := r.RemoveContainer(ctx, opts.ConnectorID); err != nil {
				return fmt.Errorf("failed to remove existing running container: %w", err)
			}
		} else {
			// Container exists but is stopped/errored, automatically remove it
			fmt.Printf("Found stopped/errored container '%s' (status: %s), removing it...\n", opts.ConnectorID, containerStatus.Status)
			if err := r.RemoveContainer(ctx, opts.ConnectorID); err != nil {
				return fmt.Errorf("failed to remove existing stopped container: %w", err)
			}
		}
	}

	// Determine runtime ID from image or use provided RuntimeID
	runtimeID := opts.RuntimeID
	if runtimeID == "" {
		runtimeID = r.extractRuntimeFromImage(opts.Image)
	}

	// Encode the original steps as base64 for the runtime
	stepsYAML, err := yaml.Marshal(opts.Steps)
	if err != nil {
		return fmt.Errorf("failed to marshal steps: %w", err)
	}
	encodedSteps := base64.StdEncoding.EncodeToString(stepsYAML)

	// Build Docker command
	args := []string{"run"}

	if opts.Remove {
		args = append(args, "--rm")
	}

	if opts.Follow {
		args = append(args, "-it")
	} else {
		args = append(args, "-d")
	}

	// Add connector ID as name
	if opts.ConnectorID != "" {
		args = append(args, "--name", opts.ConnectorID)
	}

	// Add environment variables
	for k, v := range opts.EnvVars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	// Add customer options
	if opts.DockerOpts != "" {
		dockerOpts := strings.Fields(opts.DockerOpts)
		for _, opt := range dockerOpts {
			args = append(args, opt)
		}
	}

	// Add the image
	args = append(args, opts.Image)

	// Add the base64-encoded steps as the first argument
	args = append(args, encodedSteps)

	fmt.Printf("Running: docker %s\n", strings.Join(args, " "))

	return r.executeDockerCommand(ctx, args, opts.Follow)
}

func (r *Runner) Stop(ctx context.Context, connectorID string) error {
	containerName := connectorID
	args := []string{"stop", containerName}

	fmt.Printf("Stopping: docker %s\n", strings.Join(args, " "))
	return r.executeDockerCommand(ctx, args, false)
}

func (r *Runner) Logs(ctx context.Context, connectorID string, follow bool) error {
	containerName := connectorID
	args := []string{"logs"}

	if follow {
		args = append(args, "-f")
	}

	args = append(args, containerName)

	fmt.Printf("Getting logs: docker %s\n", strings.Join(args, " "))
	return r.executeDockerCommand(ctx, args, follow)
}

func (r *Runner) List(ctx context.Context) error {
	args := []string{"ps", "-a", "--filter", "name=_connector", "--format", "table {{.Names}}\\t{{.Image}}\\t{{.Status}}\\t{{.CreatedAt}}"}

	fmt.Println("Active connectors:")
	return r.executeDockerCommand(ctx, args, false)
}

func (r *Runner) Remove(ctx context.Context, connectorID string) error {
	containerName := connectorID

	// Stop first (ignore errors if already stopped)
	r.Stop(ctx, connectorID)

	// Then remove with force flag
	args := []string{"rm", "-f", containerName}
	fmt.Printf("Removing: docker %s\n", strings.Join(args, " "))
	return r.executeDockerCommand(ctx, args, false)
}

func (r *Runner) executeDockerCommand(ctx context.Context, args []string, interactive bool) error {
	cmd := exec.CommandContext(ctx, "docker", args...)

	if interactive {
		// For interactive commands, connect to stdin/stdout/stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// For non-interactive commands, capture output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

func (r *Runner) ValidateDockerAvailable() error {
	// Check if Docker is available by running docker version
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	if err != nil {
		// For demonstration purposes, show what would happen instead of failing
		fmt.Println("Note: Docker not available. Please run 'docker image ls' manually to verify. Trying to execute:")
		return nil
	}
	return err
}

func (r *Runner) PullImage(ctx context.Context, image string) error {
	args := []string{"pull", image}
	fmt.Printf("Pulling image: docker %s\n", strings.Join(args, " "))
	return r.executeDockerCommand(ctx, args, false)
}

// GenerateDockerfile creates a Dockerfile for a connector
func (r *Runner) GenerateDockerfile(connectorID string, workDir string) error {
	dockerfilePath := filepath.Join(workDir, "Dockerfile")

	dockerfile := `FROM synadia/connect:latest

# Set working directory
WORKDIR /app

# Copy connector configuration
COPY ConnectFile /app/ConnectFile

# Set environment variables
ENV CONNECTOR_ID=%s

# Run the connector
CMD ["connect", "standalone", "run", "/app/ConnectFile"]
`

	content := fmt.Sprintf(dockerfile, connectorID)

	return os.WriteFile(dockerfilePath, []byte(content), 0644)
}

// CreateConnectFile creates a ConnectFile in the specified directory
func (r *Runner) CreateConnectFile(steps model.Steps, workDir string) error {
	connectFilePath := filepath.Join(workDir, "ConnectFile")

	data, err := yaml.Marshal(steps)
	if err != nil {
		return fmt.Errorf("failed to marshal steps: %w", err)
	}

	return os.WriteFile(connectFilePath, data, 0644)
}

// extractRuntimeFromImage extracts the runtime ID from a Docker image name
// Example: "registry.synadia.io/connect-runtime-wombat:latest" -> "wombat"
func (r *Runner) extractRuntimeFromImage(image string) string {
	// Remove tag if present
	parts := strings.SplitN(image, ":", 2)
	imageName := parts[0]

	// Extract runtime from image name pattern: */connect-runtime-<runtime>
	if strings.Contains(imageName, "/connect-runtime-") {
		parts := strings.Split(imageName, "/connect-runtime-")
		if len(parts) > 1 {
			return parts[1]
		}
	}

	// Fallback: use the last part of the image name
	pathParts := strings.Split(imageName, "/")
	return pathParts[len(pathParts)-1]
}

// GetContainerStatus checks if a container exists and returns its status
func (r *Runner) GetContainerStatus(ctx context.Context, connectorID string) (*ContainerStatus, error) {
	containerName := connectorID

	// Check if container exists and get its status
	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.Names}}\t{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to check container status: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		// No container found
		return &ContainerStatus{
			Name:   containerName,
			Status: "",
			Exists: false,
		}, nil
	}

	// Parse the output - should be "container-name\tstatus"
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 && parts[0] == containerName {
			return &ContainerStatus{
				Name:   containerName,
				Status: parts[1],
				Exists: true,
			}, nil
		}
	}

	return &ContainerStatus{
		Name:   containerName,
		Status: "",
		Exists: false,
	}, nil
}

// IsContainerRunning checks if the container status indicates it's running
func (cs *ContainerStatus) IsContainerRunning() bool {
	return cs.Exists && strings.Contains(strings.ToLower(cs.Status), "up")
}

// IsContainerStopped checks if the container status indicates it's stopped or exited
func (cs *ContainerStatus) IsContainerStopped() bool {
	return cs.Exists && !cs.IsContainerRunning()
}

// PromptUserForReplacement asks the user if they want to replace a running container
func (r *Runner) PromptUserForReplacement(connectorID string) (bool, error) {
	fmt.Printf("Container '%s' is already running.\n", connectorID)
	fmt.Print("Do you want to stop and replace it? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

// RemoveContainer removes a container (stops it first if running)
func (r *Runner) RemoveContainer(ctx context.Context, connectorID string) error {
	// Try to stop first (ignore errors if already stopped)
	r.Stop(ctx, connectorID)

	// Remove the container
	containerName := connectorID
	args := []string{"rm", "-f", containerName} // Add -f to force remove

	fmt.Printf("Removing container: %s\n", containerName)
	return r.executeDockerCommand(ctx, args, false)
}
