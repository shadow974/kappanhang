package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/getopt"
)

var (
	verboseLog                bool
	quietLog                  bool
	connectAddress            string
	username                  string
	password                  string
	civAddress                byte
	controllerAddress         byte
	serialTCPPort             uint16
	enableSerialDevice        bool
	rigctldPort               uint16
	runCmd                    string
	runCmdOnSerialPortCreated string
	statusLogInterval         time.Duration
	setDataModeOnTx           bool
	debugPackets              bool
)

func parseArgs() {
	h := getopt.BoolLong("help", 'h', "display help")
	v := getopt.BoolLong("verbose", 'v', "Enable verbose (debug) logging")
	q := getopt.BoolLong("quiet", 'q', "Disable logging")
	a := getopt.StringLong("address", 'a', "IC-705", "Connect to address")
	u := getopt.StringLong("username", 'u', "beer", "Username")
	p := getopt.StringLong("password", 'p', "beerbeer", "Password")
	c := getopt.StringLong("civ-address", 'c', "0xa4", "CI-V address for radio")
	t := getopt.Uint16Long("serial-tcp-port", 't', 4531, "Expose radio's serial port on this TCP port")
	s := getopt.BoolLong("enable-serial-device", 's', "Expose radio's serial port as a virtual serial port")
	r := getopt.Uint16Long("rigctld-port", 'r', 4532, "Use this TCP port for the internal rigctld")
	e := getopt.StringLong("exec", 'e', "", "Exec cmd when connected")
	o := getopt.StringLong("exec-serial", 'o', "socat /tmp/kappanhang-IC-705.pty /tmp/vmware.pty", "Exec cmd when virtual serial port is created, set to - to disable")
	i := getopt.Uint16Long("log-interval", 'i', 150, "Status bar/log interval in milliseconds")
	d := getopt.BoolLong("set-data-tx", 'd', "Automatically enable data mode on TX")
	dp := getopt.BoolLong("debug-packets", 'D', "Show CI-V packets for debugging")
	ca := getopt.StringLong("controller-address", 'z', "0xe0", "Controller address")

	getopt.Parse()

	if *h || *a == "" || (*q && *v) {
		fmt.Println(getAboutStr())
		getopt.Usage()
		os.Exit(1)
	}

	verboseLog = *v
	quietLog = *q
	connectAddress = *a
	username = *u
	password = *p

	*c = strings.Replace(*c, "0x", "", -1)
	*c = strings.Replace(*c, "0X", "", -1)
	civAddressInt, err := strconv.ParseInt(*c, 16, 64)
	if err != nil {
		fmt.Println("invalid CI-V address: can't parse", *c)
		os.Exit(1)
	}
	civAddress = byte(civAddressInt)

	*ca = strings.Replace(*ca, "0x", "", -1)
	*ca = strings.Replace(*ca, "0X", "", -1)
	controllerAddressInt, err := strconv.ParseInt(*ca, 16, 64)
	if err != nil {
		fmt.Println("invalid CI-V address for controller: can't parse", *ca)
		os.Exit(1)
	}
	controllerAddress = byte(controllerAddressInt)

	serialTCPPort = *t
	enableSerialDevice = *s
	rigctldPort = *r
	runCmd = *e
	runCmdOnSerialPortCreated = *o
	statusLogInterval = time.Duration(*i) * time.Millisecond
	setDataModeOnTx = *d
	debugPackets = *dp
}
