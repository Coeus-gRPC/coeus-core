# Coeus

[![Release](https://img.shields.io/github/v/release/Coeus-gRPC/coeus-core?include_prereleases)](https://github.com/Coeus-gRPC/coeus-core/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/Coeus-gRPC/coeus-core)](https://goreportcard.com/report/github.com/Coeus-gRPC/coeus-core)
[![License](https://img.shields.io/github/license/Coeus-gRPC/coeus-core)](https://raw.githubusercontent.com/Coeus-gRPC/coeus-core/LICENSE.md)

Coeus, whose name means "query" and "intelligence," is a command-line tool that focuses on [gRPC](https://grpc.io/) benchmarking and load-testing.

## Features
🔍 **gRPC Focused** - Building a gRPC endpoint? Coeus is the tool to test that!

🔍 **Modern CLI Tool** - Modern command-line interface built using Go and Cobra.

🔍 **Effortless Configuration** - Build your test configuration and message using JSON.

🔍 **macOS Client** - Not a fan of Command-line? Don't worry! There is a [macOS App](https://github.com/Coeus-gRPC/coeus-macos), just like Postman.

## Install

### Download
Please check out the [release page](https://github.com/Coeus-gRPC/coeus-core/releases) and download the latest binary release.

### Compile
**Build using go**
```sh
cd coeus-core
go build .
```

## Usage
```
usage: coeus-core [<flags>]

Flags:
  -c, --config                   path to the coeus config file
  -h, --help                     Show help for root command.
Args:
  [<host>]  Host and port to test.
```
### Example
```
./coeus-core --config ./test/testdata/config/sample_config.json
```

## Configuration file
Here's how to write the configuration file. You can also find the example file under ```./test/testdata/config/sample_config.json```

Please note: all fields are required.
```yaml
{
  "id": "01f395d4-9bd32dc7aed9",                // Random ID to identify the configuration
  "totalCallNum": 10,                           // Total number of calls to make
  "concurrent": 1,                              // Concurrent call number
  "targetHost": "api.coeustool.dev:443",        // Target host, please include port address, as well
  "insecure": false,                            // If set True, Coeus-core will retrieve TLS cert from destination server
  "timeout": 1000,                              // If set -1, timeout will be ignored
  "protoFile": "./test/testdata/proto/greeter.proto", // Path to protobuf files, currently only support Protobuf syntax version 3
  "methodName": "greeterservice.Greeter.SayHello",    // It is formatted as "package_name.service_name.method_name"
  "messageDataFile": "./test/testdata/message/sample_message.json", // Format the message file as standard JSON
  "outputFilePath": "./output/output.json"            // The output file should exist (or be created) beforehand
}
```

## Sample Output
```json
{
  "reportID": "0bb89e31-0653-4fd2-bfc2-2c45d8e7e1e6",
  "totalCallNum": 500,
  "successCallCount": 500,
  "concurrencyLevel": 10,
  "totalTimeConsumption": 4636.218,
  "averageTimeConsumption": 91.34727799999999,
  "fastestTimeConsumption": 82.109,
  "slowestTimeConsumption": 108.491,
  "timeConsumptions": [
    100.719,
    99.622,
    ------emitted------
    86.922
  ],
  "distribution": {
    "10": 85.784,
    "25": 87.976,
    "5": 84.929,
    "50": 90.605,
    "75": 94.264,
    "90": 98.352,
    "95": 100.957,
    "99": 104.449
  },
  "requestPerSecond": 107.84652199684845,
  "messages": [
    "reply:\"Hello World!\"",
    "reply:\"Hello World!\"",
    ------emitted------
    "reply:\"Hello World!\""
  ]
}
```

## Future Plan
- Add result visualization
- Add more customizable options


## License
Released under [MIT License](https://www.mit.edu/~amini/LICENSE.md)