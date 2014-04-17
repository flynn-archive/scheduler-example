# Scheduler Example

This is an example scheduler for Layer 0 of Flynn to demonstrate using the Flynn APIs to create a scheduler and some of the types of jobs you might design into your scheduler.

It currently does not demonstrate any sort of smart resource allocation or job queuing, as those would be specific to requirements of your system.

## Installing

Use the `grid` tool to schedule this scheduler. You can either build and containerize this project on a Flynn host first, or let it pull from the Docker index. 

	$ grid schedule flynn/scheduler-example

It should give you a job ID and the address it's listening on. 

## Using the Scheduler

This scheduler exposes an HTTP API for scheduling jobs of three fairly common types. It also provides a means of interacting with the jobs after scheduling just to demonstrate how the job types are different.

### Batch Jobs

Batch jobs are one-off processes that run non-interactively to do some work or produce some result, such as running a build or doing data analysis. In this example, when you schedule a batch job, you are also connected to its output so you can see whatever it prints to STDOUT.

Scheduling a batch job is done with an HTTP GET on a sub-resource of `/batch` that represents a Docker container image to use. The query parameter string is used as the command or arguments to pass to the container when it is run.

	GET /batch/<image name>?<cmd args>

If scheduled, the output of the process will be streamed back via HTTP streaming. Here is an example using the example batch job primesum, which builds up the sum of N prime numbers. We can use curl:

	$ curl 10.0.2.15:55005/batch/flynn/primesum?5000

Since `flynn/primesum` is a Docker container that simply runs the primesum command, the query string is used as arguments to the primesum program. In this case we tell it to calculate the sum of the first 5000 prime numbers. 

### Service Jobs (TODO)

Service jobs are long-running processes, usually server daemons. Unlike batch jobs, after scheduling they're expected to continue to run indefinitely. If they do stop it's assumed to be a crash and will be restarted. For this type of job, the scheduler models them RESTful resources you can create and manage. 

	POST 	/services
	GET 	/services/<jobid>
	DELETE 	/services/<jobid>

### PTY Jobs

PTY jobs are interactive jobs that expect to be attached to with a PTY. This example demonstrates functionality similar to `heroku run`. For this scheduler, as a web server with no client, we expose this via a page that you load in a browser that will give you a browser-based terminal emulator to interact with the job. Otherwise, it's similar to the batch job.

	GET /pty/<image name>?<cmd args>

So for example, to open your browser to a terminal that will connect to an interactive job such as vi:

	$ open http://10.0.2.15:55005/pty/flynn/busybox?vi

## License

BSD

## Flynn 

[Flynn](https://flynn.io) is a modular, open source Platform as a Service (PaaS). 

If you're new to Flynn, start [here](https://github.com/flynn/flynn).

### Status

Flynn is in active development and **currently unsuitable for production** use. 

Users are encouraged to experiment with Flynn but should assume there are stability, security, and performance weaknesses throughout the project. This warning will be removed when Flynn is ready for production use.

Please report bugs as issues on the appropriate repository. If you have a general question or don't know which repo to use, report them [here](https://github.com/flynn/flynn/issues).

## Contributing

We welcome and encourage community contributions to Flynn.

Since the project is still unstable, there are specific priorities for development. Pull requests that do not address these priorities will not be accepted until Flynn is production ready.

Please familiarize yourself with the [Contribution Guidelines](https://flynn.io/docs/contributing) and [Project Roadmap](https://flynn.io/docs/roadmap) before contributing.

There are many ways to help Flynn besides contributing code:

 - Fix bugs or file issues
 - Improve the [documentation](https://github.com/flynn/flynn.io) including this website
 - [Contribute](https://flynn.io/#sponsor) financially to support core development

Flynn is a [trademark](https://flynn.io/docs/trademark-guidelines) of Apollic Software, LLC.
