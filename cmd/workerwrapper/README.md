WorkerWrapper
=============

The workerwrapper is a program that handles the Gearman logic for worker jobs so the
workers can be Gearman agnostic.

Interface
---------

The workerwrapper expects the worker programs it runs to implement the following interface:

Input

 - A single base64 encoded payload. This corresponds to the Gearman payload.

Output

 - The worker's "response" is the stdout of the process. This corresponds to the Gearman data field.
 - The worker's "warnings" are the last X lines of the stderr of the process. This corresponds to the Gearman warnings field.
 - The success / failure of the worker is a function of the exit code of the process.
 - Logs should be written to stderr.


Usage
-----

First build the WorkerWrapper and add it to your GOPATH with the command:
`go get github.com/Clever/baseworker-go/cmd/workerwrapper`

Make sure you have exported the Gearman environment variables:
`export GEARMAN_HOST=localhost`
`export GEARMAN_PORT=4730`

Then run it:
`workerwrapper --name job_name --cmd job_cmd`

job_name is the name of the Gearman job to listen for
job_cmd is the name of the worker program to run

And submit jobs to it:
`gearman -f job_name -h localhost -v -s`
