# Contributing

This document provides guidelines for contributing to this repository.

## Development Environment Setup

### Prerequisites
* Basic tools: `curl`, `jq`, `oc`
* go: v1.19
* Container engine: `podman`

## Running a test environment

### Kwok Container
"Kwok" is a Kubernetes SIGs-hosted project. KWOK is an abbreviation for
Kubernetes Without Kubelet. Kwok simply simulates the node's behaviour.
As a result, it can mimic a high number of nodes and pods while consuming
only a small amount of memory.

Run the following command for a fresh clone to initialize and update the submodule:
   ```
   git submodule update --init
   ```

To run the Kwok container with the Kwok Kubernetes tool, follow these steps:

1. Build the kwok container using the following command:
   ```
   podman build -t kwok -f kwok/Dockerfile kwok
   ```
2. Bring the clusters up by running the following command from the 
   repo's root directory:
    ```
    podman kube play kwok/kwok_container_default.yml
    ```

3. Check `podman pod list` should list the below pod
    ```
    POD ID        NAME            STATUS      CREATED        INFRA ID      # OF CONTAINERS
    e815836efc86  kwok-pod        Running     1 minutes ago  8466696a9956  2
    ```

4. Once the Kwok clusters are up and running, set the cluster details in the
   OpenShift client with the following commands:
    ```
    oc config set-cluster kwok-host --server=http://127.0.0.1:8080
    oc config set-cluster kwok-m01 --server=http://127.0.0.1:8070
    oc config set-cluster kwok-rh01 --server=http://127.0.0.1:8060
    ```

5. Create new contexts (you only need to set the contexts once) for the Kwok
   clusters with the following commands:
    ```
    oc config set-context kwok-host --cluster=kwok-host
    oc config set-context kwok-m01 --cluster=kwok-m01
    oc config set-context kwok-rh01 --cluster=kwok-rh01
    ```

6. Set the Kwok context as the current context, if you've previously switched
   to another cluster, with the following command:
    ```
    oc config use-context { kwok-host, kwok-m01, kwok-rh01 }
    ```

Now you can access the cluster using kubectl, e.g.: `kubectl get ns`.

> **Note**
>
> Alternatively, you can also use the kubeconfig file on the repository providing the
> context and the cluster inline with the command. For example:
>
> ```oc --kubeconfig=./kwok/kubeconfig --context=kwok-host get ns```

### Setting up a containerized Splunk instance

To set up a containerized Splunk instance, you can use either podman or docker.
Follow these steps:

1. Build the Splunk container using the following command:
   ```
   podman build -t splunk ./splunk
   ```
2. Run the Splunk instance by running the below command:
   ```
   podman play kube splunk/splunk_container_default.yaml
   ```
   For more information about the specific command inputs and options, refer
   to the [Splunk documentation][CS1].
3. Once the container is running, you can log in to the Splunk instance using
   the username `admin` and the password you set up.
4. To access the Splunk [REST API][CS2],
   you can make API calls from outside the container using the `curl` command.
   For example, you can run the following command to search all data:
     ```
     curl -u admin:YourPassword -k https://localhost:8089/services/search/jobs -d search="search *"
     ```
    you may also use the dedicated .netrc file and authenticate with it instead:
    ```
     curl --netrc-file -k https://localhost:8089/services/search/jobs -d search="search *"
    ```
5. If you want to use the Splunk UI, open a web browser on the host and navigate to
   `localhost:8000`.

[CS1]:
https://docs.splunk.com/Documentation/Splunk/9.0.4/Installation/DeployandrunSplunkEnterpriseinsideDockercontainers
[CS2]:
https://docs.splunk.com/Documentation/Splunk/9.0.4/RESTTUT/RESTTutorialIntro

### Building and running the segment-bridge container image

The scripts in this repo can be built into a container image to enable
scheduling and running them on K8s clusters.

To build the image locally, one needs to be logged in to a `redhat.com` account
(With e.g `podman login`) in order to access the base image and then the image
can be built with:
```
podman build -t segment-bridge .
```

The scripts require access to Splunk, Segment and OpenShift credentials. One
way to provide such access is to mount the local `~/.netrc` and
`~/.kube/config` files (Assuming they contain suitable credentials) to the
image with a command like the following:
```
podman run -it --rm \
         -v ~/.netrc:/usr/local/etc/netrc:z \
         -v ~/.kube/config:/usr/local/etc/kube_config:z \
         segment-bridge
```
The following command can be run inside the container to test the full chain of
scripts. This will copy real data from the staging audit logs in Splunk into the
Segment DEV environment:
```
fetch-uj-records.sh | splunk-to-segment.sh | segment-mass-uploader.sh
```

### Unit Tests
Go unit tests are included in various packages within the repository.
Go unit tests are located within the tests directory, with filenames ending with
_tests.go and with .go.

#### Running the Unit Tests Locally
1. Clone your fork of the project.
2. Initialise [git-submodules](#Running-a-test-environment) before running the unit tests. 
3. Navigate to the project's root directory
4. To run all the Go unit tests in the repository,
execute the following command `go clean -testcache && go test ./...`
a similar output is expected:

    ```
    ?   	github.com/redhat-appstudio/segment-bridge.git/cmd/querygen	[no test files]
    ok  	github.com/redhat-appstudio/segment-bridge.git/querygen	0.002s
    ok  	github.com/redhat-appstudio/segment-bridge.git/queryprint 0.002s
    ok  	github.com/redhat-appstudio/segment-bridge.git/scripts	0.002s
    ok  	github.com/redhat-appstudio/segment-bridge.git/segment	0.002s
    ok  	github.com/redhat-appstudio/segment-bridge.git/stats	0.002s
    ok  	github.com/redhat-appstudio/segment-bridge.git/webfixture	0.002s
    ```
4. If you want to run tests for a specific directory/path, you can do so by providing
the package path, like this:
    ```
    go test ./querygen
    ```

#### Test Coverage
[TBD]

### Integration Tests
[TBD]

#### Running the Integration Tests
[TBD]

### Before submitting the PR

1. The repository enforces pre-commit checks. Ensure installing `pre-commit` using `pip install -r requirements.lock` running `pre-commit run --all-files` and fixing any issues raised before committing any changes.
2. Ensure to run `gofmt` to format your code.
3. Make sure all unit tests are passing.

### Commit Messages
We use [gitlint](https://jorisroovers.com/gitlint/) to standardize commit messages,
following the [Conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) format.

If you include a Jira ticket identifier (e.g., RHTAPWATCH-387) in the commit message,
PR name, or branch name, it will link to the Jira ticket.

```
feat(RHTAPWATCH-387): Include the UserAgent field

Include the UserAgent field in all events sent to Segment.
Also a small fix for `get-workspace-map.sh` to improve local usage.

Signed-off-by: Your Name <your-email@example.com>

```

### Pull Request Description
When creating a Pull Request (PR), use the commit message as a starting point,
and add a brief explanation. Include what changes you made and why.
This helps reviewers understand your work without needing to investigate
deeply. Clear information leads to a smoother review process.

### Code Review Guidelines
* Each PR should be approved by at least 2 team members. Those approvals are only
relevant if given since the last major change in the PR content.

* All comments raised during code review should be addressed (fixed/replied).
  * Reviewers should resolve the comments that they've raised once they think
    they were properly addressed.
  * If a comment was addressed by the PR author but the reviewer did not resolve or
    reply within 1 workday (reviewer's workday), then the comment can be resolved by
    the PR author or by another reviewer.

* All new and existing automated tests should pass.

* A PR should be open for at least 1 workday at all time zones within the team. i.e.
team members from all time zones should have an opportunity to review the PR within
their working hours.

* When reviewing a PR, verify that the PR addresses these points:
  * Edge cases
  * Race conditions
  * All new functionality is covered by unit tests
  * It should not be necessary to manually run the code to see if a certain part works,
    a test should cover it
  * The commits should be atomic, meaning that if we revert it, we don't lose something
    important that we didn't intend to lose
  * PRs should have a specific focus. If it can be divided into smaller standalone
    PRs, then it needs to be split up. The smaller the better
  * Check that the added functionality is not already possible with an existing
    part of the code
  * The code is maintainable and testable
  * The code and tests do not introduce instability to the testing framework
