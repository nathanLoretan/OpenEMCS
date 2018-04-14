package main

import(
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// Log file, flush the file every 30 minutes
var logfile *os.File
var tickerLog *time.Ticker = time.NewTicker(time.Minute * time.Duration(30))
var log_filePath string 		= "log.txt"

// Define the process for measurement
const defaultInterval = 0
var tickerMeasurement *time.Ticker  
var stop chan bool = make(chan bool, 1) // Stop tickerMeasurement

// Define struct for json encoding
var Io controls
var n node

var serverIp_filePath string 	= "serverIp.txt"

var clientNode = &http.Client{
	Timeout: time.Second * 10,
}

func main(){	
	// var err error	// USE ONLY FOR DEFINE LOGFILE
		
	fmt.Println("-----------------------------------------")
	fmt.Println("---          Node openEMCS            ---")
	fmt.Println("-----------------------------------------")
	
	// Open Logfile
	/*os.Remove(log_filePath)
	logfile, err = os.OpenFile(log_filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file, %s\n", err)
	}
	log.SetOutput(logfile)
	go FlushLog()*/
		
	// initialize node and node's server
	DatabaseInit()
	ConfigInit()
	GPIOInit()
	GetSreverIP()
	
	go ServerInit()
	go TakeMeasurements()
	
	Connect()
	
	var cmd string	
	// Infinity Loop, Stop if the user enter 'exit'
	for ; cmd != "exit"; {
		fmt.Print(">>> ")
		fmt.Scanln(&cmd)
	}
	
	Disconnect();
	tickerMeasurement.Stop()
	tickerLog.Stop()
	logfile.Close()
}

/*															*
* Function/Interface: GetSreverIP
* Param:					
* Return:					
* Description: 	search the server's ip address saved into the	
				file serverIp.txt
*															*/
func GetSreverIP()(error) {
	var err error
	
	buffer, err := ioutil.ReadFile(serverIp_filePath)
	if(err != nil) {
		log.Println(err)
		return err
	}
	
	// Clear the string
	serverIp = stripchars(string(buffer), "\t\n\r\x00 abcdefghijclmnopqrstuvwxyz")
	return err
	
}

/*															*
* Function/Interface: connect
* Param:					
* Return:					
* Description: 	send the nodes IP to the server that save it
				into the database
*															*/
func Connect()(){
	
	var status int
	for ; status != http.StatusOK; {
		
		// Encode the body for request
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(n)		
		
		// Send request and get response
		req, _ := http.NewRequest(post, resourceHttp + serverIp + serverPort + resourcesNodes, body)
		resp, err := clientNode.Do(req)	
		
		// If the node is connected
		if(err == nil) {
			json.NewDecoder(resp.Body).Decode(&n)
			DatabaseUpdateNode(n)
		
			// Close network connection
			resp.Body.Close()		
			status = resp.StatusCode
			
		} else {
			status = http.StatusInternalServerError 
		}
		
	}
	fmt.Println("\n===> Node connected to the server!!!\n")
}

/*															*
* Function/Interface: Disconnect
* Param:					
* Return:					
* Description: 	delete the node saving into te server's database
*															*/
func Disconnect()(){
	
	var status int
	for ; status != http.StatusOK; {
		// Send request and get response
		req, _ := http.NewRequest(delete, resourceHttp + serverIp + serverPort + resourcesNodes, nil)
		resp, err := clientNode.Do(req)
		
		// If the node is deconnected
		if(err == nil) {
			// Close network connection
			resp.Body.Close()
			status = resp.StatusCode
			
		} else {
			status = http.StatusInternalServerError 
		}
		
	}	
	fmt.Println("\n===> Node disconnected to the server!!!")
}

/*															*
* Function/Interface: GPIOInit
* Param:					
* Return:					
* Description: 	Initialize pin's config. of the node, all 
				pins are disable, no description and value
				null.
*															*/
func GPIOInit()(){
	var err error

	GPIOSetup()
	Io, err = DatabaseGetControls()
	GPIOConfig(Io)
	
	if err != nil {
		log.Println(err)
	}
}

/*															*
* Function/Interface: GPIOConfig
* Param:			
	- IoBuffer		the new configurations for the IO
* Return:					
* Description:	Configate the node's IO with new config.
				sent by the user. Write the value on 
				digital output and DAC.
*															*/
func GPIOConfig(IoBuffer controls)(){

	copy(Io.ADConvert[:], IoBuffer.ADConvert[:])
	copy(Io.DAConvert[:], IoBuffer.DAConvert[:])
	copy(Io.DigitalIO[:], IoBuffer.DigitalIO[:])
	copy(Io.DigitalOut[:], IoBuffer.DigitalOut[:])
	
	for i := range Io.DAConvert {
		if(Io.DAConvert[i].Enable) {
			if(Io.DAConvert[i].Reference > 0) {
				WriteDA(DAChannel[i], int(Io.DAConvert[i].Value * convertResolution / Io.DAConvert[i].Reference))
			} else {
				WriteDA(DAChannel[i], 0)
			}
		}
	}
	
	for i := range Io.DigitalOut {
		if(Io.DigitalOut[i].Enable) {
			WritePin(OutPin[i], Io.DigitalOut[i].Value)
		}
	}
	
	for i := range Io.DigitalIO {
		if(Io.DigitalIO[i].Enable && Io.DigitalIO[i].Mode == out) {
			WritePin(IOPin[i], Io.DigitalIO[i].Value)
		}
	}
	
	err := DatabaseUpdateControls(Io);
	if err != nil {
		log.Println(err)
	}
}

/*															*
* Function/Interface: GPIOUpdate
* Param:					
* Return:					
* Description:	Read the value of digital input and ADC
*															*/
func GPIOUpdate()(){

	for i := range Io.ADConvert {
		if(Io.ADConvert[i].Enable) {
			Io.ADConvert[i].Value = float32(ReadAD(ADChannel[i])) / convertResolution * Io.ADConvert[i].Reference
		}
	}
	
	for i := range Io.DigitalIO {
		if(Io.DigitalIO[i].Enable && Io.DigitalIO[i].Mode == in) {
			Io.DigitalIO[i].Value = ReadPin(IOPin[i])
		}
	}
}

/*															*
* Function/Interface: ConfigInit
* Param:					
* Return:					
* Description: 	Initialize configuration of the node: name,
				description, interval between measurements 	
*															*/
func ConfigInit()(){
	var err error
	
	n, err = DatabaseGetNode()
	if err != nil {
		log.Println(err)
	}
	
	if(n.Config.Interval > 0){
		SetTickerMeasurement(n.Config.Interval)
		stop <- false
	
	} else {
		if(tickerMeasurement != nil) {
			tickerMeasurement.Stop()
		}
		stop <- true
	
	}
}

/*															*
* Function/Interface: configInit()
* Param:					
* Return:					
* Description: 	change configuration of the node: name,
				description, interval between measurements
*															*/
func ConfigureNode(c configuration)(){

	n.Config.Name = c.Name
	n.Config.Description = c.Description
	n.Config.Interval = c.Interval
	
	err := DatabaseUpdateNode(n);
	if err != nil {
		log.Println(err)
	}
	
	if(n.Config.Interval > 0){
		SetTickerMeasurement(n.Config.Interval)
		stop <- false
		
	} else {
		if(tickerMeasurement != nil) {
			tickerMeasurement.Stop()
		}
		stop <- true
	}
}

/*															*
* Function/Interface: setTickerMeasurement()
* Param:		
  m:	time in minutes
* Return:					
* Description: 	set the time for the ticker that call the	
				function TakeMeasurement. 
*															*/
func SetTickerMeasurement(s int)(){
	if(s == 0) {
		log.Println("Error: try to set ticker with time <= 0")
		return
	}
	tickerMeasurement = time.NewTicker(time.Second * time.Duration(s))
}

/*															*
* Function/Interface: TakeMeasurements()
* Param:					
* Return:					
* Description: 	take the current and voltage measurements
				and save the value into the database
*															*/
func TakeMeasurements()(){
	
	for{
		if(tickerMeasurement != nil) {
			
			select {
				case <- stop :
				
				case <- tickerMeasurement.C :
		
				var m measurement
				m.V 	= float32(ReadAD(VChannel) / convertResolution * measurementReference)
				m.Vrms	= float32(ReadAD(VrmsChannel) / convertResolution * measurementReference)
				m.I 	= float32(ReadAD(IChannel) / convertResolution * measurementReference)
				m.Irms	= float32(ReadAD(IrmsChannel) / convertResolution * measurementReference)
				m.Time	= time.Now().Unix()
				
				err := DatabaseAddMeasurements(m)

				if err != nil {
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

// from https://www.rosettacode.org/wiki/Strip_a_set_of_characters_from_a_string#Go
func stripchars(str, chr string) string {
    return strings.Map(func(r rune) rune {
        if strings.IndexRune(chr, r) < 0 {
            return r
        }
        return -1
    }, str)
}
