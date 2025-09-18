/*
========================
 Graceful Shutdown Notes
========================

OS Signals Overview:
--------------------
1. os.Interrupt
   - High-level Go constant.
   - Usually triggered by pressing Ctrl+C in the terminal.
   - Internally maps to SIGINT on most systems.

2. syscall.SIGINT
   - Low-level constant for the interrupt signal.
   - Sent when the user presses Ctrl+C.
   - Equivalent to os.Interrupt.

3. syscall.SIGTERM
   - Termination signal.
   - Sent by `kill <pid>` or by Docker/Kubernetes when stopping a container.
   - Used to request a graceful shutdown (cleanup before exit).
   - If the process does not exit in time, the OS may send SIGKILL (force quit).

Best Practices in Production:
-----------------------------
- Always handle both SIGINT (local dev, Ctrl+C) and SIGTERM (prod/Docker/K8s).
- On receiving a signal:
    • Stop accepting new requests.
    • Allow in-flight requests to finish.
    • Close DB connections, caches, message queues, etc.
    • Flush logs or metrics if needed.
- Implement a safety timeout (e.g., 10–30s) to force exit if cleanup hangs.
- In Kubernetes:
    • SIGTERM is sent before pod termination.
    • Use readiness/liveness probes to avoid traffic during shutdown.
    • Ensure cleanup logic runs before container is killed.

Signal Handling Flow:
---------------------
1. Process receives SIGINT (Ctrl+C) or SIGTERM (K8s/Docker stop).
2. App stops accepting new work.
3. Cleanup tasks run (DB close, cache release, etc.).
4. Process exits cleanly with code 0.

Summary:
--------
- SIGINT (Ctrl+C) → manual interruption, common in dev.
- SIGTERM → graceful termination in production (Docker/K8s).
- SIGKILL → force quit, no cleanup possible (avoid relying on this).
- Handle both SIGINT and SIGTERM for consistent behavior across dev and prod.
*/

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shubham1-Kumar/students-apis/internal/config"
	"github.com/Shubham1-Kumar/students-apis/internal/http/handlers/student"
	"github.com/Shubham1-Kumar/students-apis/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialize", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux() // it returns server mux whch is your router

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.Addr))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start sever")
		}
	}()

	<-done

	// shutdown logic
	slog.Info("Shubtting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failded to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
