package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	sig "os/signal"
	"syscall"

	log "github.com/pion/ion-log"
	"github.com/pion/ion/cmd/biz/json-rpc/server"
	"github.com/spf13/viper"
)

var (
	conf server.Config
	file string
)

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {config file}")
	fmt.Println("      -h (show help info)")
}

func unmarshal(rawVal interface{}) bool {
	if err := viper.Unmarshal(rawVal); err != nil {
		fmt.Printf("config file %s loaded failed. %v\n", file, err)
		return false
	}
	return true
}

func load() bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file %s read failed. %v\n", file, err)
		return false
	}

	if !unmarshal(&conf) || !unmarshal(&conf.Config) {
		return false
	}
	if err != nil {
		fmt.Printf("config file %s loaded failed. %v\n", file, err)
		return false
	}

	fmt.Printf("config %s load ok!\n", file)

	return true
}

func parse() bool {
	flag.StringVar(&file, "c", "conf/conf.toml", "config file")
	help := flag.Bool("h", false, "help info")
	flag.Parse()
	if !load() {
		return false
	}

	if *help {
		showHelp()
		return false
	}
	return true
}

func main() {
	if !parse() {
		showHelp()
		os.Exit(-1)
	}

	fixByFile := []string{"asm_amd64.s", "proc.go"}
	fixByFunc := []string{}
	log.Init(conf.Log.Level, fixByFile, fixByFunc)

	log.Infof("--- starting biz node ---")

	s := server.NewServer(conf)
	if err := s.Start(); err != nil {
		log.Errorf("biz start error: %v", err)
		os.Exit(-1)
	}
	defer s.Close()

	// Press Ctrl+C to exit the process
	ch := make(chan os.Signal, 1)
	sig.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}
