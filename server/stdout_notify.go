package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Line single line text
type Line struct {
	Stream string `json:"stream"`
}

// StdoutNotify notify to STDOUT
type StdoutNotify struct {
}

// SendBuildReport stream build output to STDOUT
func (out StdoutNotify) SendBuildReport(ctx context.Context, r io.ReadCloser, target BuildTarget) {
	defer r.Close()
	// create build report
	var buildReport BuildReport
	buildReport.BuildTarget = target
	// print build status
	fmt.Println("===== Docker Build =====")
	fmt.Printf("Building %s:%s\n", target.Name, target.Tag)
	fmt.Printf("From git context: %s\n", target.GitContext)
	buildReport.Start = time.Now()
	// stream build output
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		var line Line
		if err := json.Unmarshal([]byte(s), &line); err != nil {
			log.Error(err)
			break
		}
		fmt.Print(line.Stream)
	}
	buildReport.Duration = time.Since(buildReport.Start)
	// TODO: decide on build status
	buildReport.Status = "Completed"
	// send build report stats
	gStats.SendReport(buildReport)
	// print duration
	fmt.Printf("Build duration: %s\n", buildReport.Duration)
	if err := scanner.Err(); err != nil {
		log.Error(err)
	}
}

// SendPushReport print push details
func (out StdoutNotify) SendPushReport(ctx context.Context, r io.ReadCloser, image string) {
	defer r.Close()

	// print push status
	fmt.Println("===== Docker Push =====")
	fmt.Printf("Pushing %s ...\n", image)

	start := time.Now()
	// stream build output
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		var line Line
		if err := json.Unmarshal([]byte(s), &line); err != nil {
			log.Error(err)
			break
		}
		fmt.Print(line.Stream)
	}

	// print duration
	fmt.Printf("Push duration: %s\n", time.Since(start))
	if err := scanner.Err(); err != nil {
		log.Error(err)
	}
}
