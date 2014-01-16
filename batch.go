package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	services "github.com/flynn/go-discoverd"
	"github.com/flynn/go-dockerclient"
	hosts "github.com/flynn/lorne/client"
	"github.com/flynn/lorne/types"
	jobs "github.com/flynn/sampi/client"
	"github.com/flynn/sampi/types"
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
	sched, err := jobs.New()
	if error500(w, err) {
		return
	}

	s, err := services.Services("flynn-lorne", services.DefaultTimeout)
	if error500(w, err) {
		return
	}
	hostid := s[0].Attrs["id"]
	jobid := jobs.RandomJobID("batch")

	host, err := hosts.New(hostid)
	if error500(w, err) {
		return
	}

	conn, attachWait, err := host.Attach(&lorne.AttachReq{
		JobID: jobid,
		Flags: lorne.AttachFlagStdout | lorne.AttachFlagStderr | lorne.AttachFlagStdin | lorne.AttachFlagStream,
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

	schedReq := &sampi.ScheduleReq{
		Incremental: true,
		HostJobs: map[string][]*sampi.Job{
			hostid: {{ID: jobid, Config: &config}}},
	}
	_, err = sched.Schedule(schedReq)
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
