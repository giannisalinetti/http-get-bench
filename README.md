# http-get-bench: parallel HTTP GET bencharmks made easy

This is a simple tool to perform parallel HTTP GET test
and record response times for every call and for the overall test.

## Build and install

To build and install http-get-bench:

```
$ make
$ sudo make install
```

This will run tests and compile the binary install it under */usr/local/bin* on your system.

## Usage

At least an URL must be provided to begin the benchmark. The following example
will run a series of 10 parallel HTTP GET calls on **www.example.com**.

```
$ http-get-bench --url www.example.com -n 10
```

The expected output will be similar to this:

```
Performing host DNS lookup...	Done
Beginning benchmark...

| Attempt: 6	| Url: http://www.example.com	| Status: 200 OK	| Time: 238172358 |
| Attempt: 1	| Url: http://www.example.com	| Status: 200 OK	| Time: 242034046 |
| Attempt: 7	| Url: http://www.example.com	| Status: 200 OK	| Time: 242244306 |
| Attempt: 9	| Url: http://www.example.com	| Status: 200 OK	| Time: 246435924 |
| Attempt: 3	| Url: http://www.example.com	| Status: 200 OK	| Time: 250223883 |
| Attempt: 2	| Url: http://www.example.com	| Status: 200 OK	| Time: 250141894 |
| Attempt: 4	| Url: http://www.example.com	| Status: 200 OK	| Time: 250140194 |
| Attempt: 5	| Url: http://www.example.com	| Status: 200 OK	| Time: 254471200 |
| Attempt: 10	| Url: http://www.example.com	| Status: 200 OK	| Time: 254644109 |

Benchmark completed.
Total number of requests:			10
Total elapsed time in nanoseconds:		254687651
```

The tool will beging with a DNS lookup test on the provided hostname. After the test the 
goroutines running the http.Get() methods will be spawned and the resultds recorded into a 
buffered channel.
The order of attempts varies upon the response time.


If needed, it is possibile to print the html output with the **-p** flag:

```
$  http-get-bench --url www.example.com -n 10 -p
```



