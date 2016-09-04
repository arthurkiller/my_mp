package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/arthurkiller/my_mp/db"
	"github.com/arthurkiller/my_mp/grpc/news"
	"github.com/arthurkiller/my_mp/grpc/profile"
	"google.golang.org/grpc"
)

func main() {
	//gen the config
	redism.RedismConf{
		Maxactive:   500,
		Maxidle:     300,
		Idletimeout: 60,
		Masters:     []string{"127.0.0.1:6379"},
		Slaves:      [][]string{{"127.0.0.1:6379", "127.0.0.1:6379"}},
	}

	//makeup the redis connetion pool manager

	//startup the server
	grpc.NewServer()
	newsServer := new(news.NewsServer())
	news.RegisterNewsServer()
	profile.RegisterProfileServer()

	//catch the sys signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			<-sigc
		}
	}()

	return 0
}
