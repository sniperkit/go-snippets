
[![Docker Stars](https://img.shields.io/docker/stars/swce/metadata-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/swce/metadata-resource/stars/count/)
[![Docker pulls](https://img.shields.io/docker/pulls/swce/metadata-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/swce/metadata-resource)
[![Docker build status](https://img.shields.io/docker/build/swce/metadata-resource.svg)](https://github.com/swce/metadata-resource)
[![Docker Automated build](https://img.shields.io/docker/automated/swce/metadata-resource.svg)](https://github.com/swce/metadata-resource)

[![dockeri.co](http://dockeri.co/image/swce/metadata-resource)](https://hub.docker.com/r/swce/metadata-resource/)

# Concourse CI Metadata Resource

Implements a resource that passes to a task the metadata of the job.

Caution: misuse may result in angry concourse developers. This resource was created for the sole purpose of linking the Artifactory artifacts to the current build.

Opinionated pipeline suggestion [here](#opinionated-pipeline)

## Thanks

This resource was implemented based on the [build-metadata-resource](https://github.com/vito/build-metadata-resource)

## Source Configuration

``` YAML
resource_types:
  - name: meta
    type: docker-image
    source:
      repository: swce/metadata-resource
      
resources:
  - name: meta
    type: meta
```

#### Parameters

*None.*

## Behavior

### `check`: Produce a single dummy key

Produce the current timestamp to invalidate the previous version so every build will get a fresh and relevant copy of the metadata.

### `in`: Write the metadata to the destination dir

 - "$BUILD_ID" > build-id
 - "$BUILD_NAME" > build-name
 - "$BUILD_JOB_NAME" > build-job-name
 - "$BUILD_PIPELINE_NAME" > build-pipeline-name
 - "$BUILD_TEAM_NAME" > build-team-name
 - "$ATC_EXTERNAL_URL" > atc-external-url 

#### Parameters

*None.*

### `out`: Unsed

Unused

#### Parameters

*None.*

## Examples

```YAML
resource_types:
  - name: meta
    type: docker-image
    source:
      repository: swce/metadata-resource

resources:
  - name: meta
    type: meta

jobs:

  - name: build
    plan:
      - get: meta
      - task: build
        file: tools/tasks/build/task.yml


```

The build job gets in the `meta` dir all the files with the respected values in them to use as it pleases

## Opinionated pipeline

Use this resource to link the artifacts created by the build step to the current build. This is helpfull in a couple of ways: 
 - Artifactory will show the build number in the metadata of the artifacts, which helps understanding which build created the artifacts.
 - When working with snapshots and production artifactory repositories, we can easily promote the artifacts of the build from snapshot to production without the need to realize all the artifacts created by the build and copying them one by one. This is done using the promote api call of Artifactory.
 
We will use the [keyval-resource](https://github.com/swce/keyval-resource) to pass the build number to the step that will release the artifact to production.

The `pipeline.yml` file:

```YAML
resource_types:
  - name: meta
    type: docker-image
    source:
      repository: swce/metadata-resource
  - name: keyval
    type: docker-image
    source:
      repository: swce/keyval-resource

resources:
  - name: meta
    type: meta
  - name: keyval
    type: keyval

jobs:

  - name: build
    plan:
      - get: meta
      - task: build
        file: tools/tasks/build/task.yml
      - put: keyval
        params:
          file: keyvalout/keyval.properties

...

  - name: prod-deploy
    plan:
      - get: keyval
      - task: prod-deploy
        file: tools/tasks/prod-deploy/task.yml

```

The `build` task:

```sh

pipeline_id=`cat "${ROOT_FOLDER}/meta/build-name"`
echo "Pipeline id is $pipeline_id"
export "PASSED_PIPELINE_ID=$pipeline_id"

...

gradlew ... -PbuildId="${PASSED_PIPELINE_ID}"

```

The `build.gradle` file of the project:
```gradle

artifactory {
  publish {
    publishBuildInfo = true
  }
  clientConfig.info.setBuildNumber(buildId)
  clientConfig.publisher.addMatrixParam(BuildInfoFields.BUILD_NUMBER, buildId)
}

```

The `prod-deploy` task:

```sh
  ...
        echo "Promoting build to production repo"
        local appName=$(retrieveAppName)
        local args="{\"status\": \"Deployed\",\"comment\": \"moving to production\",\"copy\": true,\"sourceRepo\": \"${REPO_SNAPSHOT}\",\"targetRepo\": \"${REPO_RELEASE}\",\"properties\": {\"retention.pinned\":[\"7\"]}}"
        curl --fail -u "${artifactory_user}:${artifactory_password}" -H "Content-Type: application/json" -X POST -d "'$args'" "${artifactory_contextUrl}/api/build/promote/${appName}/${PASSED_PIPELINE_ID}"

```


## Development

### Prerequisites

* golang is *required* - version 1.9.x is tested; earlier versions may also
  work.
* docker is *required* - version 17.06.x is tested; earlier versions may also
  work.
* godep is used for dependency management of the golang packages.

### Running the tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t metadata-resource .
```

### Contributing

Please make all pull requests to the `master` branch and ensure tests pass
locally.
