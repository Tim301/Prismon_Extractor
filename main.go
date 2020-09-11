package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var tRefresh *int
var help *bool
var Url string

func PrintHelp() {
	print(`Need 2 arguments:
Arg1 = Url to scan.
Arg2 = name/path of csv where log will be save.

Parameter:
-t int, Set refresh periode in seconde.
`)
	os.Exit(0)
}

func main() {
	tRefresh = flag.Int("t", 120, "Set refresh periode in seconde")
	help = flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help || len(os.Args) < 3 {
		PrintHelp()
	}
	Url = flag.Args()[0]
	fmt.Println("URL: " + Url)
	CSVPath := flag.Args()[1]
	fmt.Println("CSV: " + CSVPath)
	fmt.Print("Refresh each " + strconv.Itoa(*tRefresh) + " seconds")
	file, _ := os.OpenFile(CSVPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	test := GetVQA(Url)
	file.WriteString(test)
	file.Close()
	for true {
		time.Sleep(time.Duration(*tRefresh) * time.Second)
		fmt.Println("Refreshed at " + time.Now().Format("3:04PM"))
		t1 := time.Now()
		test := GetVQA(Url)
		t2 := time.Now()
		//file, _ := os.OpenFile("test.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		scannerTop := bufio.NewScanner(strings.NewReader(test))
		scannerTop.Split(bufio.ScanLines)
		cmd := exec.Command("sed", "-n", "2p", CSVPath)
		LastLine, err := cmd.CombinedOutput()
		if err != nil {
			LastLine = nil
			print("error")
		}
		fmt.Print(string(LastLine[:]))
		line := 0
		var testsplit string
		var t3 time.Time
		//scannerTop.Scan()
		for scannerTop.Scan() {
				if scannerTop.Text()[0:len(scannerTop.Text())] == string(LastLine[0:len(scannerTop.Text())]) {
					//print("youpi")
					break
				} else {
					testsplit = testsplit + scannerTop.Text() + "\n"
				}
			line++
		}
		t3 = time.Now()
		if line > 2 {
			data, _ := ioutil.ReadFile(CSVPath)
			//fmt.Print(string(data))
			testsplit = testsplit + string(data)[58:]
			file, _ := os.OpenFile(CSVPath, os.O_RDWR, 0644)
			file.WriteAt([]byte(testsplit), 0)
			file.Close()
		} else {
			print("Something seems to be wrong!?")
		}
		t4 := time.Now()
		fmt.Println("Request time: " + t2.Sub(t1).String() + ", Insert time: " + t4.Sub(t3).String() + ", Compute time: " + t4.Sub(t2).String())
	}
os.Exit(0)
}

func GetVQA(url string) (values string) {
	values = ""
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 90 * time.Second, Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		values = "Error"
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		values = bodyString
	}
	return
}
