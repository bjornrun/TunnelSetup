/*
TunnelSetup

Copyright (c) 2013 Bjorn Runaker

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/* Changes:
1.1 Added SOCKS5 setup
1.0 Initial version
*/

package main

import (
	"fmt"
	"github.com/stvp/go-toml-config"
	"flag"
	"os"
	"os/exec"
	"log"
	"io"
	"os/user"
	"bufio"
)

var (
	instance            = config.Int("instance", 0)
	proxyPort           = config.Int("proxy.port", 1080)
	proxyServerAddr     = config.String("proxy.address", "10.0.1.136")
	proxySSHMasterFlag  = config.String("proxy.sshmasterflag", "-o \"ControlMaster=yes\" -o \"ControlPath=~/.ssh/%r@%h:%p\"")

)


var cfgFile string
var command string
var ctrlSocket string
var tunnelListFile string
var bSocks bool

var Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
    flag.PrintDefaults()
}


func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func main() {
	flag.StringVar(&cfgFile, "c", "tunnels.cfg", "Tunnel config setup file")
	flag.StringVar(&command, "e", "help", "Execute command: help|attach|detatch|config|forward <local port:ip:remote port>|remote <remote port:ip:local port>")
	flag.BoolVar(&bSocks,"s", false, "Enable SOCKS server on attach")
	
	flag.Usage = Usage
    flag.Parse()

	fmt.Printf("Tunnel Setup\n")
		
	if err := config.Parse(cfgFile); err != nil {
		panic(err)
	}
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
	ctrlSocket := fmt.Sprintf("%s/.ssh/%s.%d", usr.HomeDir, *proxyServerAddr, *instance)	
	tunnelListFile := fmt.Sprintf("%s/.ssh/%s.%d.txt", usr.HomeDir, *proxyServerAddr, *instance)
	
	if command == "help" {
		flag.PrintDefaults()
		os.Exit(0)
	} else
	if command == "attach" {		
		
		if _, err := os.Stat(ctrlSocket); err == nil {
			fmt.Printf("Server %s already attached", *proxyServerAddr)
			os.Exit(1)
		}
		var cmd *exec.Cmd
		if bSocks {
			cmd = exec.Command("ssh", "-o", "ControlMaster=yes", "-o", fmt.Sprintf("ControlPath=%s", ctrlSocket),"-fNT","-D", fmt.Sprintf("%d", *proxyPort), *proxyServerAddr)				
    

		} else
		{
			cmd = exec.Command("ssh", "-o", "ControlMaster=yes", "-o", fmt.Sprintf("ControlPath=%s", ctrlSocket),"-fNT", *proxyServerAddr)	
     
			
		}
					     stdout, err := cmd.StdoutPipe()

		 if err != nil {
        log.Fatal(err)
     }
     stderr, err := cmd.StderrPipe()
     if err != nil {
        log.Fatal(err)
     }
     err = cmd.Start()
     if err != nil {
        log.Fatal(err)
     }
     

    go io.Copy(os.Stdout, stdout)
    go io.Copy(os.Stderr, stderr)

	err = cmd.Wait()
	if err != nil {
        log.Fatal(err)
     }

		fmt.Printf("Server %s is now attached\n", *proxyServerAddr)
		os.Exit(0)
	} else
	if command == "detach" {
		fmt.Printf("detach")
		if _, err := os.Stat(ctrlSocket); os.IsNotExist(err)  {
			fmt.Printf("Server %s already detached", *proxyServerAddr)
			os.Exit(1)
		}
		
		cmd := exec.Command("ssh", "-O", "stop", "-o", fmt.Sprintf("ControlPath=%s", ctrlSocket), *proxyServerAddr)

     stdout, err := cmd.StdoutPipe()
     if err != nil {
        log.Fatal(err)
     }
     stderr, err := cmd.StderrPipe()
     if err != nil {
        log.Fatal(err)
     }
     err = cmd.Start()
     if err != nil {
        log.Fatal(err)
     }
     

    go io.Copy(os.Stdout, stdout)
    go io.Copy(os.Stderr, stderr)

	err = cmd.Wait()
	if err != nil {
        log.Fatal(err)
     }
		fmt.Printf("Server %s is now detached\n", *proxyServerAddr)
		os.Remove(tunnelListFile)
		os.Exit(0)
	} else
	if command == "forward" {
//		fmt.Printf("forward tunnel\n")
		if _, err := os.Stat(ctrlSocket); os.IsNotExist(err)  {
			fmt.Printf("Server %s is not attached", *proxyServerAddr)
			os.Exit(1)
		}
		
		
		cmd := exec.Command("ssh", "-O", "forward", "-o", fmt.Sprintf("ControlPath=%s", ctrlSocket), "-L", os.Args[len(os.Args) - 1], *proxyServerAddr)
     stdout, err := cmd.StdoutPipe()
     if err != nil {
        log.Fatal(err)
     }
     stderr, err := cmd.StderrPipe()
     if err != nil {
        log.Fatal(err)
     }
     err = cmd.Start()
     if err != nil {
        log.Fatal(err)
     }
     

    go io.Copy(os.Stdout, stdout)
    go io.Copy(os.Stderr, stderr)

	err = cmd.Wait()
	if err != nil {
        log.Fatal(err)
     }
		fmt.Printf("Forward tunnel %s active\n", os.Args[len(os.Args) - 1])
	f, err := os.OpenFile(tunnelListFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, err = os.OpenFile(tunnelListFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
		    panic(err)
		}
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("Forward %s\n", os.Args[len(os.Args) - 1])); err != nil {
	    panic(err)
	}
			os.Exit(0)
	} else
	if command == "remote" {
//		fmt.Printf("remote tunnel\n")
		if _, err := os.Stat(ctrlSocket); os.IsNotExist(err)  {
			fmt.Printf("Server %s is not attached", *proxyServerAddr)
			os.Exit(1)
		}
		
		
		cmd := exec.Command("ssh", "-O", "forward", "-o", fmt.Sprintf("ControlPath=%s", ctrlSocket), "-R", os.Args[len(os.Args) - 1], *proxyServerAddr)
     stdout, err := cmd.StdoutPipe()
     if err != nil {
        log.Fatal(err)
     }
     stderr, err := cmd.StderrPipe()
     if err != nil {
        log.Fatal(err)
     }
     err = cmd.Start()
     if err != nil {
        log.Fatal(err)
     }
     

    go io.Copy(os.Stdout, stdout)
    go io.Copy(os.Stderr, stderr)

	err = cmd.Wait()
	if err != nil {
        log.Fatal(err)
     }
		fmt.Printf("Remote tunnel %s active\n", os.Args[len(os.Args) - 1])
		f, err := os.OpenFile(tunnelListFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, err = os.OpenFile(tunnelListFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
	    	panic(err)
		}
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("Remote %s\n", os.Args[len(os.Args) - 1])); err != nil {
	    panic(err)
	}
		os.Exit(0)
	} else
	if command == "config" {
		fmt.Printf("Configuration:\nInstance: %d\nServer: %s:%d\n",*instance,*proxyServerAddr,*proxyPort)
		if _, err := os.Stat(ctrlSocket); os.IsNotExist(err)  {
			fmt.Printf("Not attached\n")
			os.Exit(0)
		} else
		{
			fmt.Printf("Attached to Proxy %s\n", *proxyServerAddr)			
		}

		lines, err := readLines(tunnelListFile)
    if err == nil {
		if len(lines) > 0 {
			fmt.Printf("Tunnels:\n")
			for _,item := range lines {
        		fmt.Println(item)
			}
		}
		}
		os.Exit(0)
	} else
	{
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}	
