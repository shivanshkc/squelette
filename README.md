# Squelette

## Introduction

Squelette (french for Skeleton) is a quickstart template designed to simplify the process of creating web services in Go. Its goal is to provide developers with a well-structured project layout, essential dependencies, and configuration files.

## Getting Started

1. Clone the repository:
    ```sh
    git clone https://github.com/shivanshkc/squelette.git
    cd squelette
    ```

2. Replace `squelette` with your desired project name in all files and directories.

3. Rename the `cmd/squelette` folder to your desired project name.

4. Create a configs file by running:
    ```sh
    cp configs/configs.sample.yaml configs/configs.yaml
    ```

5. Run using:
    ```sh
    make run
    ```

## Makefile Commands

The `Makefile` includes several commands to streamline common tasks:

- `make build`: Build the project.
- `make run`: Compile and run the project.
- `make image`: Build the container image of the project.
- `make container`: Run an application container.
- `make test`: Run tests for the project.
- `make lint`: Run linters to check code quality.

## Adding an API

TODO
