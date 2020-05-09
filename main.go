package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
)

type options struct {
	Project  string `short:"p" long:"project" description:"Project"`
	Instance string `short:"i" long:"instance" description:"Instance"`
	Database string `short:"d" long:"database" description:"Database"`
}

const (
	defaultGRPCPort = 9010
	defaultRESTPort = 9020
)

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			os.Exit(0)
		}
		exitf("invalid options")
	}

	if opts.Project == "" || opts.Instance == "" || opts.Database == "" {
		exitf("missing parameters: -p, -i, -d are required")
	}

	emulator := &Emulator{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		GRPCPort: defaultGRPCPort,
		RESTPort: defaultRESTPort,
	}
	gCloudConfig := NewGCloudConfig()
	go waitInterrupt(emulator, gCloudConfig)

	fmt.Println("Start spanner emulator...")
	if err := emulator.Start(); err != nil {
		exitf("failed to start emulator: %v", err)
	}

	fmt.Println("Wait for emulator to be up...")
	if err := emulator.WaitForReady(); err != nil {
		exitf("failed to wait for emulator: %v", err)
	}

	fmt.Printf("Create spanner instance: %s\n", opts.Instance)
	if err := createInstance(defaultRESTPort, opts.Project, opts.Instance); err != nil {
		exitf("failed to create spanner instance: %v", err)
	}
	fmt.Printf("Create spanner database: %s\n", opts.Database)
	if err := createDatabase(defaultRESTPort, opts.Project, opts.Instance, opts.Database); err != nil {
		exitf("failed to create spanner database: %v", err)
	}

	fmt.Printf("Create an ephemeral gcloud configuration: %s\n", gCloudConfig.Name)
	if err := gCloudConfig.Setup(opts.Project, defaultRESTPort); err != nil {
		exitf("failed to create gcloud configuration: %v", err)
	}

	fmt.Printf(`You can use the following environment variables to access the emulator.

export CLOUDSDK_ACTIVE_CONFIG_NAME=%s
export SPANNER_EMULATOR_HOST=localhost:%d

Now emulator is ready.
`, gCloudConfig.Name, defaultGRPCPort)

	if err := emulator.WaitForFinish(); err != nil {
		exitf("failed to wait for emulator to be finished: %v", err)
	}
}

func waitInterrupt(emulator *Emulator, gCloudConfig *GCloudConfig) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Shutting down...")
	if err := emulator.Shutdown(); err != nil {
		exitf("failed to shut down emulator: %v", err)
	}

	fmt.Printf("Delete an ephemeral gcloud configuration: %s\n", gCloudConfig.Name)
	if err := gCloudConfig.CleanUp(); err != nil {
		exitf("failed to delete gcloud configuration: %v", err)
	}
}

func exitf(format string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, "ERROR: " + fmt.Sprintf(format, a...))
	os.Exit(1)
}
