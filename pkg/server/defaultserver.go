package server

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
)

const (
	defaultWorkerNum = 15
	defaultQueueSize = 1000
)

type BusinessJob interface {
	Execute()
}

var _ Server = &DefaultServer{}

type DefaultServer struct {
	container    *restful.Container
	listener     map[int]net.Listener
	logger       *zap.Logger
	logLock      *sync.RWMutex
	workerNum    int
	jobQueue     chan BusinessJob
	jobqueuesize int
}

// Container returns container of server
func (s *DefaultServer) Container() *restful.Container {
	return s.container
}

// Start starts server
func (s *DefaultServer) Start() {
	// start audit workers on the background
	if s.workerNum <= 0 {
		s.workerNum = defaultWorkerNum
	}
	if s.jobqueuesize <= 0 {
		s.jobqueuesize = defaultQueueSize
	}
	queue := make(chan BusinessJob, s.jobqueuesize)
	for i := 0; i < s.workerNum; i++ {
		go func() {
			for {
				job, ok := <-queue
				// if queue is closed, stop the worker
				if !ok {
					return
				}
				job.Execute()
			}
		}()
	}
	s.jobQueue = queue

	for port, listener := range s.listener {
		if port <= 0 {
			continue
		}
		if listener != nil {
			// TODO: give options to enable a more complex server
			// instead of using the default one
			go http.Serve(listener, s.Container())
		} else {
			go http.ListenAndServe(fmt.Sprintf(":%d", port), s.Container())
		}
	}
	select {}
}

// SetLogger sets a zap.Logger instance on server
func (s *DefaultServer) SetLogger(logg *zap.Logger) {
	s.logLock.Lock()
	defer s.logLock.Unlock()
	s.logger = logg
}

// Log returns a zap.Logger
func (s *DefaultServer) Log() *zap.Logger {
	return s.L()
}

// L returns a zap.Logger
func (s *DefaultServer) L() *zap.Logger {
	s.logLock.RLock()
	defer s.logLock.RUnlock()
	return s.logger
}

func (s *DefaultServer) Execute() {
	//TODO: do default job
}
