# lustre_exporter
exporter for lustre filesystem

## how this works
### collectorType
This is where the magic happens.
This is responsible to parsing the right data from the relevant files

### collectors
This is a simple definitions of which files need to be read and parsed. 
It holds a list of collectorType instances that point to a specific file.
It also has some params like globs to find the files, but also regex to parse labels from the path itself.
In addition, each instance is tagged with a label that allows to enable and disable collection features in the config file.

### consts
includes some constants such as base paths to counter files, and so on

### lustre_exporter.go
runs through all the collectors, running the collect function of each one, which in turn, runs the collect function of each collectorType instance.
