package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
	tnt "github.com/tarantool/go-tarantool"
)

type Options struct {
	User    string `long:"user" description:"User to connect to the database"`
	Pass    string `long:"pass" description:"Password to connect to the database"`
	Host    string `long:"host" description:"Host to connect to the database"`
	Port    int    `long:"port" description:"Port to connect to the database"`
	IpsFile string `long:"ips-file" description:"File with ips to check"`
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if opts.User == "" || opts.Pass == "" || opts.Host == "" || opts.Port == 0 || opts.IpsFile == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	fh, err := os.Open(opts.IpsFile)
	if err != nil {
		log.Fatal("Can't open ips file: ", err)
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

	count := 0
	contains := 0

	start := time.Now()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		ip := scanner.Text()
		resp, err := conn.Call17("acl_v4_contains", []interface{}{ip})
		if err != nil {
			log.Fatal("Error while checking ip: ", err)
		}
		if resp.Data[0] != "OK" {
			log.Fatal("Error while checking ip: ", resp.Data[1])
		}
		if resp.Data[1] == true {
			contains++
		}
		count++
	}
	duration := time.Since(start)

	per := float32(0)
	if count > 0 {
		per = 100 * float32(contains) / float32(count)
	}
	rps := float64(count) / float64(duration.Seconds())
	log.Printf("Checked %d ips, %d (%.1f %%) in acl, rps: %.1f", count, contains, per, rps)
}
