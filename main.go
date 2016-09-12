package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/arthurkiller/my_mp/conf"
	"github.com/arthurkiller/my_mp/db"
	"github.com/arthurkiller/my_mp/grpc/news"
	"github.com/arthurkiller/my_mp/grpc/profile"
	"github.com/arthurkiller/my_mp/news"
	"github.com/arthurkiller/my_mp/profile"
	"google.golang.org/grpc"
)

func main() {
	pprof := flag.Bool("pprof", false, "set if turn on the pprof")
	flag.Parse()

	//gen the config
	conf := conf.Config{
		NewsServerAddr:    "127.0.0.1:23588",
		ProfileServerAddr: "127.0.0.1:23589",

		UIDInfo:  []string{"127.0.0.1:16380"},
		UIDInfoS: [][]string{{"127.0.0.1:16380"}},

		NewsIDInfo:  []string{"127.0.0.1:16381"},
		NewsIDInfoS: [][]string{{"127.0.0.1:16381"}},

		MeipaiIDInfo:  []string{"127.0.0.1:16382"},
		MeipaiIDInfoS: [][]string{{"127.0.0.1:16382"}},

		UIDBox:  []string{"127.0.0.1:16383"},
		UIDBoxS: [][]string{{"127.0.0.1:16383"}},

		UIDSelfbox:  []string{"127.0.0.1:16384"},
		UIDSelfboxS: [][]string{{"127.0.0.1:16384"}},

		UIDFans:  []string{"127.0.0.1:16385"},
		UIDFansS: [][]string{{"127.0.0.1:16385"}},

		UIDFollow:  []string{"127.0.0.1:16386"},
		UIDFollowS: [][]string{{"127.0.0.1:16386"}},

		Maxactive:   500,
		Maxidle:     300,
		Idletimeout: 60,
	}

	redisConf := redism.RedismConf{
		Maxactive:   conf.Maxactive,
		Maxidle:     conf.Maxidle,
		Idletimeout: conf.Idletimeout,
		Masters:     map[string]([]string){},
		Slaves:      map[string]([][]string){},
	}

	redisConf.Masters = make(map[string]([]string), 10)
	redisConf.Slaves = make(map[string]([][]string), 10)

	redisConf.Masters["uid-box"] = conf.UIDBox
	redisConf.Masters["uid-info"] = conf.UIDInfo
	redisConf.Masters["uid-selfbox"] = conf.UIDSelfbox
	redisConf.Masters["newsid-info"] = conf.NewsIDInfo
	redisConf.Masters["uid-fans"] = conf.UIDFans
	redisConf.Masters["uid-follow"] = conf.UIDFollow
	redisConf.Masters["meipaiid-info"] = conf.MeipaiIDInfo

	redisConf.Slaves["uid-box"] = conf.UIDBoxS
	redisConf.Slaves["uid-info"] = conf.UIDInfoS
	redisConf.Slaves["uid-selfbox"] = conf.UIDSelfboxS
	redisConf.Slaves["uid-fans"] = conf.UIDFansS
	redisConf.Slaves["uid-follow"] = conf.UIDFollowS
	redisConf.Slaves["newsid-info"] = conf.NewsIDInfoS
	redisConf.Slaves["meipaiid-info"] = conf.MeipaiIDInfoS

	//makeup the redis connetion pool manager
	redispool := redism.NewRedism(redisConf)

	//set up the pprof server
	if *pprof {
		runtime.SetBlockProfileRate(1)
		lis, err := net.Listen("tcp", ":6062")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("server profiling address: ", lis.Addr())
		go func() {
			if err := http.Serve(lis, nil); err != nil {
				fmt.Println(err)
				return
			}
		}()
	}

	//startup the server
	newsconn, err := net.Listen("tcp", conf.NewsServerAddr)
	defer newsconn.Close()
	if err != nil {
		log.Println("error in listen ", conf.NewsServerAddr)
	}
	newsServer := grpc.NewServer()
	newsS := newsm.NewNewsServer(redispool)
	news.RegisterNewsServer(newsServer, newsS)
	go func() {
		newsServer.Serve(newsconn)
	}()
	log.Println("News server started!")

	profileconn, err := net.Listen("tcp", conf.ProfileServerAddr)
	defer profileconn.Close()
	if err != nil {
		log.Println("error in listen ", conf.ProfileServerAddr)
	}
	profileServer := grpc.NewServer()
	profileS := profilem.NewProfileServer(redispool)
	profile.RegisterProfileServer(profileServer, profileS)
	go func() {
		profileServer.Serve(profileconn)
	}()
	log.Println("Profile server started!")

	//catch the sys signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		<-sigc
		redispool.Close()
		log.Println("news server is stoping")
		newsconn.Close()
		log.Println("profile server is stoping")
		profileconn.Close()
		log.Println("redispoll is stoping")
		log.Println("server stoped, bye bye!")
		return
	}
}
