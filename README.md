spanemu
===
Cloud Spanner Emulator Wrapper.

This wrapper command does the following actions to help you to run the [Cloud Spanner Emulator](https://cloud.google.com/spanner/docs/emulator) with a single command.

1. Start emulator
2. Wait for the emulator to be up
3. Create Spanner instance
4. Create Spanner database
5. Create an ephemeral gcloud configuration

## Usage

```
$ spanemu -p test-project -i test-instance -d test-db
```

## Install

```
$ go get -u github.com/yfuruyama/spanemu
```

## Roadmap

* Allow --host-port and --rest-port options
* Data persistence with spanner-cli and spanner-dump.

## Disclaimer
This is not an official Google product.
