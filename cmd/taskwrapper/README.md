TaskWrapper
=============

The taskwrapper is a program that handles the Gearman logic for any task-style program we
run so that the tasks themselves can be Gearman agnostic.

Interface
---------

The workerwrapper expects the tasks it runs to implement the following interface:

Input

 - The input arguments specified in the Gearman payload. For example if a job was submitted with the following payload:
"-h localhost -p 27017 -f s3_path". It would be translated into the corresponding command line arguments.

Output

 - The worker's "response" is the stdout of the process. This corresponds to the Gearman data field.
 - The worker's "warnings" are the last X lines of the stderr of the process. This corresponds to the Gearman warnings field.
 - The success / failure of the worker is a function of the exit code of the process.
 - Logs should be written to stderr.


Usage
-----

First build the WorkerWrapper and add it to your GOPATH with the command:
`go get github.com/Clever/baseworker-go/cmd/workerwrapper`

Then run it:
`workerwrapper --name function_name --cmd function_cmd`

function_name is the name of the Gearman function to listen for
function_cmd is the name of the task to run
gearman-host and gearman-port are optional parameters.

And submit jobs to it:
`gearman -f function_name -h localhost -v -s`
