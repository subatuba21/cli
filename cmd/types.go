package cmd

type (

	// Authentication validation response
	HasuraResponse struct {
		Path  string `json:"path"`
		Error string `json:"error"`
		Code  string `json:"code"`
	}

	// Container service
	Container struct {
		Image               string
		Name                string
		Command             []string
		Entrypoint          string
		Environment         map[string]interface{}
		Ports               []string
		Restart             string
		User                string
		Volumes             []string
		DependsOn           []string `yaml:"depends_on,omitempty"`
		EnvFile             []string `yaml:"env_file,omitempty"`
		Build               map[string]string
		HealthCheckEndpoint string
	}

	/*
		// Container services
		Services struct {
			Containers map[string]Container `yaml:"services,omitempty"`
			Version    string
		}
	*/

	ExecResult struct {
		StdOut   string
		StdErr   string
		ExitCode int
	}
)
