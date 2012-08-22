package main

///import dbg "fmt"
import "flag"
import "log"
import "os"
import "time"

import "github.com/mewmew/breyting/conf"

// confPath is the path to the config file. The default is:
//    $HOME/.config/breyting/breyting.ini
var confPath string

func init() {
	defaultConfPath := os.Getenv("HOME") + "/.config/breyting/breyting.ini"
	flag.StringVar(&confPath, "conf", defaultConfPath, "Path to config file.")
	flag.Parse()
}

func main() {
	err := breyting()
	if err != nil {
		log.Fatalln(err)
	}
}

func breyting() (err error) {
	err = conf.Reload(confPath)
	if err != nil {
		return err
	}
	go conf.Watch(confPath)
	/// ### [ tmp ] ###
	///   - server.Listen() should be added here.
	/// ### [/ tmp ] ###
	time.Sleep(1 * time.Hour)
	return nil
}
