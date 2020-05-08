package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"os/signal"
)

type options struct {
	Project  string `short:"p" long:"project" description:"Project"`
	Instance string `short:"i" long:"instance" description:"Instance"`
	Database string `short:"d" long:"database" description:"Database"`
}

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		exitf("Invalid options\n")
	}

	// TODO: --host-port and --rest-port
	if opts.Project == "" || opts.Instance == "" || opts.Database == "" {
		exitf("Missing parameters: -p, -i, -d are required\n")
	}

	// Steps:
	// 1. Start emulator
	// 2. Wait for the emulator to be up
	// 3. Create instance and database
	// 4. Import data (if any)
	// 5. Create gcloud config
	// 6. Show gcloud config to the user
	// 7. Now all steps done, user can interact with the emulator

	// TODO: check if gcloud is installed

	emulator := Emulator{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	gCloudConfig := NewGCloudConfig()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Println("Shutting down...")
		// TODO: dump database

		if err := emulator.Shutdown(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Delete an ephemeral gcloud configuration: %s\n", gCloudConfig.Name)
		if err := gCloudConfig.CleanUp(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Start spanner emulator...")
	if err := emulator.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wait for emulator to be up and ready...")
	if err := emulator.WaitForReady(); err != nil {
		exitf("failed to wait for emulator to be up: %v", err)
	}

	fmt.Printf("Create Cloud Spanner Instance: %s\n", opts.Instance)
	if err := createInstance(9020, opts.Project, opts.Instance); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Create Cloud Spanner Database: %s\n", opts.Database)
	if err := createDatabase(9020, opts.Project, opts.Instance, opts.Database); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create an ephemeral gcloud configuration: %s\n", gCloudConfig.Name)
	if err := gCloudConfig.Setup(opts.Project); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(`Now emulator is ready. You can set the following environment variables to access the emulator.
# gcloud
export CLOUDSDK_ACTIVE_CONFIG_NAME=%s
# Cloud Spanner tools
export SPANNER_EMULATOR_HOST=localhost:9010
`, gCloudConfig.Name)

	// TODO: wait for shutdown
	// if emulator process is killed, this process is also gracefully shutdowned.
	if err := emulator.WaitForFinish(); err != nil {
		log.Fatal(err)
	}
}

func exitf(format string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, "ERROR: " + fmt.Sprintf(format, a...))
	os.Exit(1)
}
