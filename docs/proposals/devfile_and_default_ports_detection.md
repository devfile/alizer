# Devfile port detection & default ports

## Related Github Issues

- https://github.com/devfile/api/issues/1154

## Background
In order to perform port detection component-wise, `alizer` is performing different strategies of port detection per component. Currently, we support 3 different strategies:

* Dockerfile Port Detection
* Docker Compose Port Detection
* Source Code Port Detection

At the moment, `alizer` doesn't support port detection from `devfile.y(a)ml` files. This means that all ports defined inside a file like this will not be included in alizer's response. _As a result, the first goal of this proposal is to include the devfiles as part of the port detection algorithm_, exactly like `dockerfile` and `compose`.

On the other hand, right now `alizer` is not returning any default ports of detected frameworks. As a result, for every component using a framework that has default ports, alizer will not include them in its response even if we haven't detected any ports from our strategies (ports list is empty). That said, as a second goal of this proposal, we would like to add a new feature in alizer that will include default ports in alizer's response.

## Devfile Port Detection
For devfile port detection we should follow the same flow with dockerfile:

1. Create a function that will:
    * look for a devfile inside a given component (The `utils.GetLocations` function should be updated to support devfiles too).
    * It will also find all ports defined in `components.kubernetes.endpoints` & `components.container.endpoints`.

```go
// GetPortsFromDevfile returns a slice of port numbers from Devfile in the given directory.
func GetPortsFromDevfile(root string) []int {
	locations := utils.GetLocations(root)
	for _, location := range locations {
		…
            return utils.ReadPortsFromDevfile(file)
	}
	return []int{}
}
```

2. Create a new `const` inside `model/model.go`:
```go
const Devfile PortDetectionAlgorithm = 3
```

3. Add a new algorithm case to the following enricher DoEnrichComponent functions:
```
dotnet_enricher.go
go_enricher.go
java_enricher.go
javascript_enricher.go
php_enricher.go
python_enricher.go
```

Following those steps will add support for devfile port detection for every component that we will be able to detect.

## Default Ports
In order to support default ports for port detection we will need to complete the following steps:

1. We will need to define a new flag `--add-default-ports` inside the `component` command. This way someone can set that they would like to have a response which will include default ports.

2. For every `framework` that we detect we should investigate if it uses any ports by default. After the investigation every `detector` should have the following function:
```go
func (d LaravelDetector) GetDefaultPorts() []string {
	return []int{8080}
}
```

3. Upon `DoPortsDetection` the default ports should be added in the end of the `ports array` in alizer's response:

```json
// default ports = 3000, 3001

// Case1: Only default ports where returned
[
    {
            "Name": "example-project",
            "Path": "/path/to/example-project/",
            "Languages": [...],
            "Ports": [3000, 3001]
    }
]

// Case2: Ports 2001, 2002 are detected for example-project
[
    {
            "Name": "example-project",
            "Path": "/path/to/example-project/",
            "Languages": [...],
            "Ports": [2001, 2002, 3000, 3001]
    }
]

```

_In case the `--add-default-ports` is not set the default value will be `false`_