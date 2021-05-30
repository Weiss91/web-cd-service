# web-cd-service
This service provides endpoints to execute bazel jobs like build, run and test for given projects. With this endpoints and correct configuration you can build a deployment pipeline for projects that use bazel as build tool.

## How it works
The base image of the web-cd-service is created in [deploy_win_template.sh](deploy_win_template.sh) and adds this service on top of this. When running the shellscript, a dockerfile is created and executed. You can adjust this template that it fits your needs. Currently it was tested for golang and java builds.

The [run_template.sh](run_template.sh) contains a sample configuration that is needed for running this service. Change it for your needs. More details in the [configuration section](##Configuration)

## Endpoints
Currently there are following endpoints:
* POST (ExecAPIKey) /execute/task --> this endpoint expects in the moment some information e.g. (assuming correct git and docker configuration): 
```json
{
	"Remote": "https://github/your-remote/repo.git",
	"Commit": "c0a702ccj3b78bdabb58da416eaf2ba217kl88e6",
	"Target": "bazelwebapi:push_image",
    "BazelCmd": "run",
    "Prio": "dev",
    "Registry": "dev-image-registry"
}
```
* GET (ReadAPIKey) /getstate/task/** --> ** is the taskID that is given by execute/task 
* GET (ReadAPIKey) /get/task/** --> ** is the taskID that is given by execute/task
* GET /health/ready
* GET /health/live

All endpoints (except health) are secured by apiKeys when set in configuration. This has to be set in the header with 'x-api-key'.
## Configuration
To get web-cd-service running you have to make some configurations. It is completely configurable by environment variables. 

GIT_USER, GIT_PASSWORD, GIT_DATA_PATH
(--> the path where the git repos should be stored that bazel can execute it. Defaults to /git)

DOCKER_CONF:
```go
// DOCKER_CONF has to be a json string matching this go struct:
type dockerConf struct {
	Path       string               `json:"path"` // the path where the config.json is stored temporary while executing push_image targets.
	Registries map[string]*registry `json:"registries"` //string = Registry
}

type registry struct {
	Registry string
	User     string `json:"user"`
	Password string `json:"password"`
	Auth     string `json:"auth"` //generated out of User and Password
}
```

API_KEY_READ, API_KEY_EXEC, SERVER_PORT (defaults to 8088)
--> the Port on which the web-api is running

STORAGE_PATH
--> the directory where the active and historic tasks are stored

## Planned Features/Work
* Tests tests tests
* Dashboard that shows some views and filters of open/historic tasks
* Gitlab/Github/Gogs Webhooks instead/additionally of /execute/task
* Swagger UI for endpoints
* Make Build targets available for download not only pushing to registries like in the moment
* Send a webhook to other CI components that an image or a component is available
* helm chart to deploy this service more easy in kubernetes/openshift environments
* Expand /execut/task endpoint that it accepts zip/gzip payloads with code instead of downloading from git
