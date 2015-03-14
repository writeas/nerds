package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net"
	"os"
	"os/exec"

	"github.com/writeas/nerds/store"
)

var (
	banner    []byte
	outDir    string
	staticDir string
	debugging bool
	rsyncHost string
	db        *sql.DB
)

const (
	colBlue   = "\033[0;34m"
	colGreen  = "\033[0;32m"
	colBGreen = "\033[1;32m"
	colCyan   = "\033[0;36m"
	colBRed   = "\033[1;31m"
	colBold   = "\033[1;37m"
	noCol     = "\033[0m"

	hr = "————————————————————————————————————————————————————————————————————————————————"
)

func main() {
	// Get any arguments
	outDirPtr := flag.String("o", "/var/write", "Directory where text files will be stored.")
	staticDirPtr := flag.String("s", "./static", "Directory where required static files exist.")
	rsyncHostPtr := flag.String("h", "", "Hostname of the server to rsync saved files to.")
	portPtr := flag.Int("p", 2323, "Port to listen on.")
	debugPtr := flag.Bool("debug", false, "Enables garrulous debug logging.")
	flag.Parse()

	outDir = *outDirPtr
	staticDir = *staticDirPtr
	rsyncHost = *rsyncHostPtr
	debugging = *debugPtr

	fmt.Print("\nCONFIG:\n")
	fmt.Printf("Output directory  : %s\n", outDir)
	fmt.Printf("Static directory  : %s\n", staticDir)
	fmt.Printf("rsync host        : %s\n", rsyncHost)
	fmt.Printf("Debugging enabled : %t\n\n", debugging)

	fmt.Print("Initializing...")
	var err error
	banner, err = ioutil.ReadFile(staticDir + "/banner.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DONE")

	// Connect to database
	dbUser := os.Getenv("WA_USER")
	dbPassword := os.Getenv("WA_PASSWORD")
	dbName := os.Getenv("WA_DB")
	dbHost := os.Getenv("WA_HOST")

	if outDir == "" && (dbUser == "" || dbPassword == "" || dbName == "") {
		// Ensure parameters needed for storage (file or database) are given
		fmt.Println("Database user, password, or database name not set.")
		return
	}

	if outDir == "" {
		fmt.Print("Connecting to database...")
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4", dbUser, dbPassword, dbHost, dbName))
		if err != nil {
			fmt.Printf("\n%s\n", err)
			return
		}
		defer db.Close()
		fmt.Println("CONNECTED")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *portPtr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on localhost:%d\n", *portPtr)

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

		if debugging {
			fmt.Print(b[0:n])
			fmt.Printf("\n%d: %s\n", n, b[0:n])
		}

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

		if debugging {
			fmt.Print(b[0:n])
			fmt.Printf("\n%d: %s\n", n, b[0:n])
		}

		if checkExit(b, n) {
			friendlyId := store.GenerateFriendlyRandomString(store.FriendlyIdLen)
			editToken := store.Generate62RandomString(32)

			if outDir == "" {
				_, err := db.Exec("INSERT INTO posts (id, content, modify_token, text_appearance) VALUES (?, ?, ?, 'mono')", friendlyId, post.Bytes(), editToken)
				if err != nil {
					fmt.Printf("There was an error saving: %s\n", err)
					output(c, "Something went terribly wrong, sorry. Try again later?\n\n")
					break
				}
				output(c, fmt.Sprintf("\n%s\nPosted! View at %shttps://write.as/%s%s", hr, colBlue, friendlyId, noCol))
			} else {
				output(c, fmt.Sprintf("\n%s\nPosted to %shttp://nerds.write.as/%s%s", hr, colBlue, friendlyId, noCol))

				if rsyncHost != "" {
					output(c, "\nPosting to secure site...")
					exec.Command("rsync", "-ptgou", outDir+"/"+friendlyId, rsyncHost+":").Run()
					output(c, fmt.Sprintf("\nPosted! View at %shttps://write.as/%s%s", colBlue, friendlyId, noCol))
				}
			}

			output(c, "\nSee you later.\n\n")
			break
		}

		if err != nil || n == 0 {
			break
		}
	}
}
