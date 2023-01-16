package main

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"github.com/jessevdk/go-flags"
	tnt "github.com/tarantool/go-tarantool"
)

type Options struct {
	User         string `long:"user" description:"User to connect to the database"`
	Pass         string `long:"pass" description:"Password to connect to the database"`
	Host         string `long:"host" description:"Host to connect to the database"`
	Port         int    `long:"port" description:"Port to connect to the database"`
	NetworksFile string `long:"networks-file" description:"File with networks to load"`
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if opts.User == "" || opts.Pass == "" || opts.Host == "" || opts.Port == 0 || opts.NetworksFile == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	fh, err := os.Open(opts.NetworksFile)
	if err != nil {
		log.Fatal("Can't open networks file: ", err)
	}
	defer fh.Close()

	tntOpts := tnt.Opts{
		User: opts.User,
		Pass: opts.Pass,
	}
	conn, err := tnt.Connect(opts.Host+":"+strconv.Itoa(opts.Port), tntOpts)
	if err != nil {
		log.Fatal("Connection to tarantool failed: ", err)
	}
	defer conn.Close()

	c := 0
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		network := scanner.Text()
		resp, err := conn.Call17("acl_v4_create_from_network", []interface{}{network, true, network})
		if err != nil {
			log.Fatal("Error while creating network: ", err)
		}
		if resp.Data[0] != "OK" {
			log.Fatal("Error while creating network: ", resp.Data[1])
		}
		c++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading networks file: ", err)
	}

	log.Printf("Loaded %d networks", c)
}
