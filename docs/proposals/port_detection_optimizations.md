# Port detection optimizations

## Related Github Issues
[Alizer port detection logic improvement](https://github.com/redhat-developer/alizer/issues/250)

## Background
The current port detection logic is prone to performance issues. An example past incident can be found in the issue [alizer#208](https://github.com/redhat-developer/alizer/issues/208) where alizer spent a big amount of time on port detection.

When [alizer#208](https://github.com/redhat-developer/alizer/issues/208) was closed, we have introduced the `excluded_folders` in [language_customization.yaml](https://github.com/redhat-developer/alizer/blob/main/go/pkg/utils/langfiles/resources/languages-customization.yml) as a workaround to reduce time of execution. This solution was focusing on the component detection part. As component detection is a different analysis process than port detection we will have to create a new more optimized solution for port detection too.

This proposal is based on the idea that for port detection we need only **some** paths and not **all** paths. More detailed, upon port detection we already know the component we are detecting ports. Another idea could be to make `port-detection` configurable and of course we can also define different strategies/algorithms for port detection.

## Optimizations
The main idea of the proposal will be to introduce 2 different optimizations regarding port detection (potential child issues in an EPIC):

* [Add specific port detection paths to each detector](#detector-port-detection-paths).
* [Define no-port-detection CLI arg](#define-no-port-detection-cli-arg).

As of now we have made changes to the port detection process for Go files/projects where it ignores any testing file with the suffix `_test.go`. Additionally, the port detection process for Go projects will ignore `mocks` and `migrations` labelled directories. This has provided significant increases in runtime speeds for Go projects.

## Detector Port Detection Paths
As a first step for better performance during port detection we can define certain paths and files for port detection and each detector. So, upon port detection process the detector will only check for configured ports inside these paths. We can apply this to all of the detectors in order to have a consistent file read process for each `DoPortDetection` function.

### Current Detector Status
One thing worth mentioning is the current detector status. We have different status for each detector and we can identify some groups of detector policies:

| Policy                                              | Detector                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
|-----------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Fetches all cached paths and checks each file.      | [echo](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/go/echo_detector.go#L36), [fasthttp](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/go/fasthttp_detector.go#L36), [gin](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/go/gin_detector.go#L36), [go-runtime](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/go/go_detector.go#L36), [fiber](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/go/gofiber_detector.go#L36), [mux](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/go/mux_detector.go#L36)                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| Looks for a specific file inside the component path | [beego](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/go/gin_detector.go#LL36C30-L36C30), [jboss-eap](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/java/jboss_eap_detector.go#L33), [micronaut](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/java/micronaut_detector.go#L49), [openliberty](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/java/openliberty_detector.go#L43), [quarkus](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/java/quarkus_detector.go#L55), [spring](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/java/spring_detector.go#L49), [vert.x](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/java/vertx_detector.go#L44), [wildfly](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/java/wildfly_detector.go#L33),  [laravel](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/php/laravel_detector.go#L32),  |
| Looks for specific paths                            | [django](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/python/django_detector.go#L30), [flask](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/python/flask_detector.go)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| Port detection not implemented yet                  | [dotnet](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/dotnet/dotnet_detector.go#L54)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| Port detection not implemented yet                  | [dotnet](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/enricher/framework/dotnet/dotnet_detector.go#L54)    

*As a result the golang related frameworks are using all component paths in order to apply port detection.*

### GetPortDetectionPaths
A way to ensure consistency on doing port detection, could be to define `GetPortDetectionPaths` inside every detector. This function could provide a list of files to the `DoPortDetection`.

```golang
// It returns a list of accepted paths or specific files from a given rootPath
func (e BeegoDetector) GetPortDetectionFiles(rootPath string) []string {
	paths := []string{"app.conf"}
	return utils.ReadComponentFiles(rootPath, paths)
}

func ReadComponentFiles(rootPath string, paths []string) []string {
// This function should merge the logic of utils.readAnyApplicationFile()
}
```

*After this update, as we will already have read all files, we will have to update the go_detector.GetPortFromFilesGo in order not to perform the same process again.*

As a result all detectors will use the same process to read specific files for port detection.

### Process Overview
The final process for each detector will be:

## Define no-port-detection CLI arg
Currently, upon `component` command we have the `port-detection` arg which specifies the port-detection strategy

### Current State
On the current state as we allow only three values for this arg:

* `docker`: Alizer will look only for dockerfile/containerfile to detect ports
* `compose`: Alizer will look only to docker-compose.yaml
* `source`: Alizer will go directly to source code for port detection.

### Deprecation of port-detection
As alizer is evolving as an application analysis tool we need to intorduce different strategies for port-detection, different than the `port-detection` flow. Right now, the port detection flow is coupled with the component detection one. For example, for each component we apply dockerfile, compose and source port detection flow. In order to decouple the functionality we would like to introduce a simpler *"port-detection by component type"*.

That said, in the future, alizer should be focusing on port detection per component type. For example, for python applications it should be focusing on the source code port detection only. All devfile/dockerfile/docker-compose files should be treated as separate components.

The above information is only a summary of the changes which will be introduced after the implementation of EPIC#1154 (https://github.com/devfile/api/issues/1154). This proposal **is related to this issue only to add a deprecation warning** for the `port-detection` CLI arg.

### The no-port-detection CLI arg
The second feature for better performance will be to allow users to have component detection without port detection afterwards. This feature will make alizer faster and will allow users choosing the amount of information they want to receive per execution.

The new CLI argument will be:
```golang
componentCmd.Flags().BoolP("no-port-detection", "n", false, "Skips the execution of port detection for all detected components. As a result no ports will be returned in the response. If it doesn't exist, alizer will run the port detection for all detected components")
```
An example usage will be:
```bash
$ ./alizer component --no-port-detection <path>
```

## The library
For users using the alizer library we can add simple documentation for using the [component_recognizer.DetectComponentsInRootWithPathAndPortStartegy](https://github.com/redhat-developer/alizer/blob/main/go/pkg/apis/recognizer/component_recognizer.go#L41) with an empty port detection strategy:

```golang
emptyStrategy := []model.PortDetectionAlgorithm{}
recognizer.DetectComponentsInRootWithPathAndPortStartegy(path, emptyStrategy, , &ctx)
```