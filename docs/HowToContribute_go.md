# Contributing to IBM MQ Cloud Go SDK 

This guide offers instructions for contributing to IBM MQ Cloud Go SDK 

## Table of Contents

1. [Coding Style](#coding-style)
2. [Commit Messages](#commit-messages)
3. [Pull Requests](#pull-requests)
4. [Updating service with new API changes](#updating-service-with-new-api-changes)
5. [Running tests](#running-tests)
    - [Unit tests](#unit-tests)
    - [Integration tests](#integration-tests)


## Coding Style

The SDK adheres to the Go coding conventions documented [here](link-to-go-coding-conventions). Run the linter using:

```bash
golangci-lint run
```

## Commit Messages

Commit messages should follow the [Angular Commit Message Guidelines](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#-commit-message-guidelines). This is because our release tool - [semantic-release](https://github.com/semantic-release/semantic-release) - uses this format for determining release versions and generating changelogs. Tools such as commitizen or commitlint can be used to help contributors and enforce commit messages. Here are some examples of acceptable commit messages, along with the release type that would be done based on the commit message:

| Commit message                                          | Release type       |
| ------------------------------------------------------- | -------------------|
| fix(resource controller): fix integration test to use correct credentials | Patch Release     |
| feat(global catalog): add global-catalog service to project | Minor Feature Release |
| feat(global search): re-gen service code with new v3 API definition |
|                                                          |                    |
| BREAKING CHANGE: The global-search service has been updated to reflect version 3 of the API. | Major Breaking Release |

## Pull Requests

If you want to contribute to the repository, follow these steps:

1. **Clone the repository**
2. **Develop and test your code changes:**
   - To build/test: `go test ./...`
   - Please add one or more tests to validate your changes.
   - Make sure everything builds/tests cleanly
3. **Commit your changes**
4. **Push to your fork and submit a pull request to the main branch**

## Updating service with new API changes

This section will guide you through the steps to generate the Go code for a service and add the generated code to the SDK.

1. **Validate the API definition** - Before processing the API definition with the SDK generator, it's recommended to validate it with the IBM OpenAPI Validator:

    ```bash
    lint-openapi -s example-service.yaml
    ```

    This command will display a list of errors and warnings found in the API definition.

2. **Run the SDK generator** - Process your API definition and generate the service and unit test code:

    ```bash
    cd <project-root>
    openapi-sdkgen.sh generate -g ibm-go -i my-service.json -o . --api-package <module-import-path>
    ```

## Running tests

The tests within the SDK consist of both unit tests and integration tests.

### Unit tests

Unit tests exercise the SDK function with local "mock" service endpoints.

- To run all the unit tests contained in the project:

    ```bash
    go test ./...
    ```

- To run a unit test for a specific package within the SDK project:

    ```bash
    cd <package-dir> (e.g. cd globalsearchv2)
    go test
    ```

### Integration tests

Integration tests use actual service endpoints deployed in IBM Cloud and therefore require the appropriate credentials.

- To generate the integration test:

    ```bash
    openapi-sdkgen.sh generate -g ibm-go -i my-service.json -o my-sdk --genITs
    ```

- To run all integration tests contained in the project:

    ```bash
    go test ./... -tags=integration
    ```

- To run the integration test for a single service:

    ```bash
    cd <package-dir> (e.g. cd globalsearchv2)
    go test -tags=integration
    ```
