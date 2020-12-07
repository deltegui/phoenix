package phoenix

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CreateStaticServer() {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
}

func PrintLogo(logoFile) {
	logo, err := ioutil.ReadFile(logoFile)
	if err != nil {
		log.Fatalf("Cannot read logo file: %s\n", err)
	}
	fmt.Println(string(logo))
}

func WaitAndStopServer(server *http.Server) {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Print("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		app.config.onStop()
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Phoenix shutdown failed:%+v", err)
	}

	log.Print("Phoenix exited properly")
}
