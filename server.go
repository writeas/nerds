package main

import (
	"fmt"
	"net"
	"bytes"
	"io/ioutil"
	"crypto/rand"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	banner []byte
)

const (
	colBlue = "\033[0;34m"
	colGreen = "\033[0;32m"
	colBGreen = "\033[1;32m"
	colCyan = "\033[0;36m"
	colBRed = "\033[1;31m"
	colBold = "\033[1;37m"
	noCol = "\033[0m"
	
	nameLen = 12
	outDir = "/var/write/"
	bannerDir = "./"
	hr = "————————————————————————————————————————————————————————————————————————————————"
)

func main() {
	ln, err := net.Listen("tcp", ":9727")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on localhost:9727")
	
	fmt.Print("Initializing...")
	banner, err = ioutil.ReadFile(bannerDir + "/banner.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DONE")
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		go handleConnection(conn)
	}
}

func output(c net.Conn, m string) bool {
	_, err := c.Write([]byte(m))
	if err != nil {
		c.Close()
		return false
	}
	return true
}

func outputBytes(c net.Conn, m []byte) bool {
	_, err := c.Write(m)
	if err != nil {
		c.Close()
		return false
	}
	return true
}

func handleConnection(c net.Conn) {
	outputBytes(c, banner)
	output(c, fmt.Sprintf("\n%sWelcome to write.as!%s\n", colBGreen, noCol))
	output(c, fmt.Sprintf("If this is freaking you out, you can get notified of the %sbrowser-based%s launch\ninstead at https://write.as.\n\n", colBold, noCol))
	
	waitForEnter(c)
	
	c.Close()
	
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
}

func waitForEnter(c net.Conn) {
	b := make([]byte, 4)
	
	output(c, fmt.Sprintf("%sPress Enter to continue...%s\n", colBRed, noCol))
	for {
		n, err := c.Read(b)
		if bytes.IndexRune(b[0:n], '\n') > -1 {
			break
		}
		if err != nil || n == 0 {
			c.Close()
			break
		}
	}
	
	output(c, fmt.Sprintf("Enter anything you like.\nPress %sCtrl-D%s to publish and quit.\n%s\n", colBold, noCol, hr))
	readInput(c)
}

func checkExit(b []byte, n int) bool {
	return n > 0 && bytes.IndexRune(b[0:n], '\n') == -1
}

func readInput(c net.Conn) {
	defer c.Close()
	
	b := make([]byte, 4096)
	
	var post bytes.Buffer
	
	for {
		n, err := c.Read(b)
		post.Write(b[0:n])
		
		if checkExit(b, n) {
			file, err := savePost(post.Bytes())
			if err != nil {
				fmt.Printf("There was an error saving: %s\n", err)
				output(c, "Something went terribly wrong, sorry. Try again later?\n\n")
				break
			}
			output(c, fmt.Sprintf("\n%s\nPosted to %shttp://nerds.write.as/%s%s\nPosting to secure site...", hr, colBlue, file, noCol))
			exec.Command("rsync", "-ptgou", outDir + file, "www:").Run()
			output(c, fmt.Sprintf("\nPosted! View at %shttps://write.as/%s%s\nSee you later.\n\n", colBlue, file, noCol))
			break
		}
		
		if err != nil || n == 0 {
			break
		}
	}
}

func savePost(post []byte) (string, error) {
	filename := generateFileName()
	f, err := os.Create(outDir + filename)
	
	defer f.Close()
	
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.WriteString(f, stripCtlAndExtFromUTF8(string(post)))
	
	return filename, err
}

func generateFileName() string {
	c := nameLen
	var dictionary string = "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, c)
	rand.Read(bytes)
	for k, v := range bytes {
		 bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r == 10 || r == 13 || (r >= 32 && r < 255) {
			return r
		}
		return -1
	}, str)
}
