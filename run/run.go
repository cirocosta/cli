package cmdrun

import (
	"errors"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/launchpad-project/cli/verbose"
	"github.com/spf13/cobra"
)

// RunCmd runs the Launchpad structure for development locally
var RunCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run Launchpad infrastructure for development locally",
	Run:     initRun,
	Example: `launchpad init`,
}

var ErrEmptyHost = errors.New("Host not found")

var detach bool

func getFromDockerHost() (string, error) {
	var dockerHost, hasDockerHost = os.LookupEnv("DOCKER_HOST")

	if !hasDockerHost {
		return "", ErrEmptyHost
	}

	var u, err = url.Parse(dockerHost)

	if err != nil {
		return "", err
	}

	host, _, err := net.SplitHostPort(u.Host)

	return host, err
}

func dockerAddressToArg() string {
	var fromDockerHost, err = getFromDockerHost()

	if err == nil {
		return fromDockerHost
	}

	if err != ErrEmptyHost {
		println("$DOCKER_HOST set with weird value. Ignoring.")
	}

	// try another way to identify if docker is running natively
	// it is very slow to do this lookup when it is not (timeout)
	_, err = net.LookupHost("docker.local")

	if err == nil {
		return "docker.local"
	}

	return "localhost"
}

func dockerAddressToArgAdd(arguments []string) {
	arguments = append(arguments, dockerAddressToArg())
}

func initRun(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		println("This command doesn't take arguments.")
		os.Exit(1)
	}

	var arguments = []string{
		"run",
		"-p", "80:80",
		"-p", "9300:9300",
		"-p", "5701:5701",
		"-p", "8001:8001",
		"-p", "8080:8080",
		"-p", "5005:5005",
	}

	var certPath, hasCertPath = os.LookupEnv("DOCKER_CERT_PATH")

	if hasCertPath {
		arguments = append(arguments, "-v")
		arguments = append(arguments, certPath+":/certs")
	}

	dockerAddressToArgAdd(arguments)

	if detach {
		arguments = append(arguments, "--detach")
	}

	arguments = append(arguments, "launchpad/dev")

	verbose.Debug("Running docker", strings.Join(arguments, " "))

	var docker = exec.Command("docker", arguments...)

	docker.Stderr = os.Stderr
	docker.Stdout = os.Stdout

	if err := docker.Run(); err != nil {
		panic(err)
	}
}

func init() {
	RunCmd.Flags().BoolVarP(&detach, "detach", "d", false,
		"Run Launchpad in background and print container ID")
}
