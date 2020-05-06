package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
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

	// TODO: check if gcloud is installed
	cmd := exec.Command("gcloud", "beta", "emulators", "spanner", "start")
	// https://stackoverflow.com/questions/33165530/prevent-ctrlc-from-interrupting-exec-command-in-golang
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Println("interrupted!!!")

		// before shutdown
		// TODO: dump database
		cmd.Process.Signal(os.Interrupt)
	}()

	fmt.Println("Start spanner emulator...")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	cloudSDKConfig, err := setupCloudSDK(opts.Project)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("export CLOUDSDK_ACTIVE_CONFIG_NAME=%s\n", cloudSDKConfig)

	if err := createInstance(9020, opts.Project, opts.Instance); err != nil {
		log.Fatal(err)
	}
	if err := createDatabase(9020, opts.Project, opts.Instance, opts.Database); err != nil {
		log.Fatal(err)
	}

	fmt.Print("export SPANNER_EMULATOR_HOST=localhost:9010\n")

	// TODO: wait for shutdown
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func exitf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
