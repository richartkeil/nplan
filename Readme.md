# NPlan

Transforms nmap XML into intermediate JSON and generates a basic network plan in the DrawIO XML format.

## Installation

Just run

```sh
go install github.com/richartkeil/nplan@latest
```

By default the binary will be installed to `$GOPATH/bin` (which you could add to your $PATH).

```sh
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

### Building the model

You can execute nplan multiple times with different `nmap` .xml files and `scan6` .txt files. The tool will initially create a .json model file that is gradually updated with the new data from each execution.

```sh
$ nplan -nmap scan1.xml
$ nplan -nmap scan2.xml -scan6 scan6.txt
$ nplan -nmap scan3.xml
```

You can generate a .drawio file when you gathered enough data with

```sh
$ nplan -export
```

### CLI Options

```sh
$ nplan -h
Usage of nplan:
  -nmap string
    	Set the path to the nmap input .xml file.
  -scan6 string
    	Set the path to the scan6 input .txt file. For this to take effect the current model should already include MAC addresses.
  -export
    	Export the current model to a .drawio file.
  -json string
    	Set the path where to store the .json model file. (default "./dist/model.json")
  -drawio string
    	Set the path were to store the exported .drawio file. (default "./dist/plan.drawio")
  -fresh
    	Delete the previous .json model and build a new one. Use with caution.
```
