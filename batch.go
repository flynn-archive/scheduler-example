package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/flynn/flynn-host/types"
	"github.com/flynn/go-dockerclient"
	"github.com/flynn/go-flynn/cluster"
)

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

func FlushWriter(writer io.Writer) flushWriter {
	fw := flushWriter{w: writer}
	if f, ok := writer.(http.Flusher); ok {
		fw.f = f
	}
	return fw
}

func batchHandler(w http.ResponseWriter, r *http.Request) {
	client, err := cluster.NewClient()
	if error500(w, err) {
		return
	}

	hosts, err := client.ListHosts()
	if error500(w, err) {
		return
	}
	var hostid string
	for id, _ := range hosts {
		hostid = id
		break
	}

	jobid := cluster.RandomJobID(strings.TrimPrefix(r.URL.Path, "/batch/"))

	hostClient, err := client.ConnectHost(hostid)
	if error500(w, err) {
		return
	}

	conn, attachWait, err := hostClient.Attach(&host.AttachReq{
		JobID: jobid,
		Flags: host.AttachFlagStdout | host.AttachFlagStderr | host.AttachFlagStdin | host.AttachFlagStream,
	}, true)
	if error500(w, err) {
		return
	}

	config := docker.Config{
		Image:        strings.TrimPrefix(r.URL.Path, "/batch/"),
		Cmd:          strings.Split(r.URL.RawQuery, "+"),
		Tty:          false,
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		OpenStdin:    false,
		StdinOnce:    false,
	}
	
	jobReq := &host.AddJobsReq{
		Incremental: true,
		HostJobs: map[string][]*host.Job{
			hostid: {{ID: jobid, Config: &config, TCPPorts: 0}}},
	}

	_, err = client.AddJobs(jobReq)
	if error500(w, err) {
		return
	}

	err = attachWait()
	if error500(w, err) {
		return
	}

	fw := FlushWriter(w)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Fprintln(fw, scanner.Text())
	}
	conn.Close()
	fw.Write([]byte("\n"))
}
