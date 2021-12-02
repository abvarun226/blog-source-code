package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abvarun226/background-workers/workers"
	"github.com/apex/log"
	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()
	graceperiod := 5 * time.Second
	workerCount := 10
	buffer := 100
	httpAddr := ":8000"

	log.Info("starting workers")
	w := workers.New(workerCount, buffer)
	w.Start(ctx)

	h := handler{worker: w}

	router := mux.NewRouter()
	router.HandleFunc("/queue-task", h.queueTask).Methods("POST")

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("starting http server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatalf("listen failed")
		}
	}()

	<-done
	log.Info("http server stopped")

	ctxTimeout, cancel := context.WithTimeout(ctx, graceperiod)
	defer func() {
		w.Stop()
		cancel()
	}()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.WithError(err).Fatalf("http server shutdown failed")
	}
}

type handler struct {
	worker workers.WorkerIface
}

func (h *handler) queueTask(w http.ResponseWriter, r *http.Request) {
	var input queueTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.WithError(err).Info("failed to read POST body")
		renderResponse(w, http.StatusBadRequest, `{"error": "failed to read POST body"}`)
		return
	}
	defer r.Body.Close()

	// parse the work duration from the request body.
	workDuration, errParse := time.ParseDuration(input.WorkDuration)
	if errParse != nil {
		log.WithError(errParse).Info("failed to parse work duration in request")
		renderResponse(w, http.StatusBadRequest, `{"error": "failed to parse work duration in request"}`)
		return
	}

	// queue the task in background task manager.
	if err := h.worker.QueueTask(input.TaskID, workDuration); err != nil {
		log.WithError(err).Info("failed to queue task")
		if err == workers.ErrWorkerBusy {
			w.Header().Set("Retry-After", "60")
			renderResponse(w, http.StatusServiceUnavailable, `{"error": "workers are busy, try again later"}`)
			return
		}
		renderResponse(w, http.StatusInternalServerError, `{"error": "failed to queue task"}`)
		return
	}

	renderResponse(w, http.StatusAccepted, `{"status": "task queued successfully"}`)
}

func renderResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

type queueTaskInput struct {
	TaskID       string `json:"task_id"`
	WorkDuration string `json:"work_duration"`
}
