# Introduction

Welcome to the Docute docs! These docs provide an in-depth explanation for everything Docute is capable of. Docute was built to replicate functionality of existing doc platforms like GitBook and mdbook but with a focus on looking good and being easy to read.

## Installation
Before installing, please make sure Go is installed. [Install Go here](https://go.dev/doc/install).

To install Docute, we can use Go's install command:
```shell
go install github.com/JackalLabs/docute@latest
```

After Docute is installed, you can run the following commands to generate your documentation. Make sure you're in the root folder of your docs before running Docute.

```shell
docute generate
```

If you want to view these generated docs, run:
```shell
docute host
```
Then you will be able to visit your docs at http://localhost:9797.

If you wish to enable live reload while editing your docs, please use:
```shell
docute watch
```
This will watch your docs for changes and refresh the page every time a change is made.