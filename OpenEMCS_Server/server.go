package main

// Compile : go run server.go switch.go handler.go resource.go database.go

import (
	"os"
	"log"
	"fmt"
	"time"
	"bufio"
	"strings"
	"net/http"
)

// Log file, flush the file every 30 minutes
var logfile *os.File
var tickerLog *time.Ticker = time.NewTicker(time.Minute * time.Duration(30)) 
var log_filePath string 		= "log.txt"

// Ping all node every minutes
var tickerPing = time.NewTicker(time.Second * 60)

func main(){
	// var err error	// USE ONLY FOR DEFINE LOGFILE
	
	fmt.Println("-----------------------------------------")
	fmt.Println("---         Server openEMCS           ---")
	fmt.Println("-----------------------------------------")
	
	// Open Logfile
	/*logfile, err = os.OpenFile(log_filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file, %s\n", err)
	}
	log.SetOutput(logfile)
	go FlushLog() */
	
	// initialize server's handler
	DatabaseInit()
	
	go ServerInit()
	go PingNode()
	
	var execute bool = true	
	// Infinity Loop, Stop if the user enter 'exit'
	for ; execute; {
		
		fmt.Printf(">>> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() 
		buffer := scanner.Text()
		cmd := strings.Fields(buffer)
		
		if(len(cmd) > 0) {
			if(cmd[0] == "exit") {
				execute = false
				
			} else {
				fmt.Printf("Error: Invalid Command %s\n", strings.Join(cmd[0:len(cmd)], " "))
			}	
		}
	}
	tickerPing.Stop()	
	tickerLog.Stop()
	logfile.Close()
}

/*															*
* Function/Interface: 	PingNode
* Param:													 
* Return:													
* Description:	Test every minutes if each nodes are connected.
				If the node is not connected, it will be removed
				from the database.
*															*/
func PingNode()(){
	
	for{
		for _ = range tickerPing.C{
			
			nodes, err := DatabaseGetAllNodes()
			if err != nil{
				log.Println(err)
			}
			
			for i := range nodes {
				ip, err :=  DatabaseGetNodeIP(nodes[i].ID)
				
				if ip != "" && err == nil {
					// Send request and get response
					req, _ := http.NewRequest(get, resourceHttp + ip + nodePort + resourceNodesPing, nil)
					resp, err := clientNode.Do(req)
					
					// if no response, delete the  node from the database
					if(err != nil || resp.StatusCode != http.StatusOK){	
						DatabaseDeleteNode(ip)
						
					// else close the connection
					}else{ 
						resp.Body.Close()		
					}					
				} else {
					log.Println(err)
				}
			}		
		}
	}	
}

/*															*
* Function/Interface: 	FlushLog
* Param:													 
* Return:													
* Description:	every 30 minutes the log file is cleared
*															*/
func FlushLog()() {
	var err error
	
	for{		
		select {
			case <- tickerLog.C :
	
			logfile.Close()
			
			// Delete file
			os.Remove(log_filePath)
			
			// Create a new file
			logfile, err = os.OpenFile(log_filePath, os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatalln("Failed to open log file, %s\n", err)
			}
			log.SetOutput(logfile)	
		}
	}
}