Cobra command CLI that contains a core for data ingest (csv/json) or gather (http)

Takes as input a yml file (which can be in the future changed/implemented of a DLA - Yacc) to provide the rules and flow the core must follow.

The core gets the data in input, does something to it (like transformation, if the rules require) then exports it by default to a Google Spreadsheet

The core allows for an extensible package system that probably uses gRPC to get the data from the package and push it to the sheet

Also, the core allows for a custom package (called storage) that takes the data from any input (including a package) and export it to the storage cloud system. This can be self-hosted or cloud-hosted.

Runs as Dockerfile (on k3, k8s, or systemd) via Cloud Run, Cloud Run Function, Labmda, etc, since this is OSS and can be self-hosted.