// Node handler

package main

import (
	"log"
	"fmt"
	"time"
	"strings"
	"strconv"
	"net/http"	
	"encoding/json"
)

// Define client and server for send/Receive request
var serverNode = http.Server{
	Addr: nodePort,
	Handler: &requestHandler{},
	ReadTimeout:  time.Second * 10,
	WriteTimeout: time.Second * 10,
}
	
//--	HANDLER		--//
//  resource: /X.X/measurments
func HandlerMeasurements(w http.ResponseWriter, r *http.Request)(){	
	
	if(r.Method == get){
		
		// get the number of measurements in parameter
		i := strings.Index(r.URL.String(), "?")
		parameters := strings.Split(r.URL.String()[i+1:], "&")
		param := strings.Split(parameters[0], "=")
		
		if(strings.Compare(param[0], "all") == 0) {
	
			nbrMeasurements, err := DatabaseGetNbrMeasurements()

			if(err == nil) {
				// configure response OK
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "application/json; charset=UTF-8")
				json.NewEncoder(w).Encode(nbrMeasurements)
				
				log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK) + " " + http.StatusText(http.StatusOK))
				fmt.Print(">>> ")
			} else {
				// configure response Interval Server Error
				fmt.Fprintf(w, responseServerError, http.StatusInternalServerError,
							http.StatusText(http.StatusInternalServerError), r.Method + " " + r.URL.String())
					
				log.Println(err)
				log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusInternalServerError) + " " + http.StatusText(http.StatusInternalServerError))
				fmt.Print(">>> ")
			} 
		
		// Get some measurements
		} else if(strings.Compare(param[0], "nbrMeasurements") == 0) {
			
			nbrMeasurements, _ := strconv.Atoi(param[1])
			// get some measurements
			measurements, err := DatabaseGetMeasurements(nbrMeasurements)

			if(err == nil) {
				// configure response OK
				w.WriteHeader(http.StatusOK)
				w.Header().Set("content-type", "application/json; charset=UTF-8")
				json.NewEncoder(w).Encode(measurements)
				
				log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK) + " " + http.StatusText(http.StatusOK))
				fmt.Print(">>> ")
			} else {
				// configure response Interval Server Error
				fmt.Fprintf(w, responseServerError, http.StatusInternalServerError,
							http.StatusText(http.StatusInternalServerError), r.Method + " " + r.URL.String())
					
				log.Println(err)
				log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusInternalServerError) + " " + http.StatusText(http.StatusInternalServerError))
				fmt.Print(">>> ")
			} 
		}  else {	
			// configure response Bad Request
			fmt.Fprintf(w, responseNodeError, http.StatusBadRequest,
					http.StatusText(http.StatusBadRequest), r.Method + " " + r.URL.String())
					
			log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusBadRequest) + " " + http.StatusText(http.StatusBadRequest))
			fmt.Print(">>> ")
		} 
							
	} else if(r.Method == delete){
		// delete measurements table		
		err := DatabaseDeleteMeasurements()
		
		if(err == nil) {
			// configure response OK
			w.WriteHeader(http.StatusOK)
			
			log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK) + " " + http.StatusText(http.StatusOK))
			fmt.Print(">>> ")
			
		} else {
			// configure response Interval Server Error
			fmt.Fprintf(w, responseServerError, http.StatusInternalServerError,
						http.StatusText(http.StatusInternalServerError), r.Method + " " + r.URL.String())
						
			log.Println(err)
			log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusInternalServerError) + " " + http.StatusText(http.StatusInternalServerError))
			fmt.Print(">>> ")
		}
		
	} else {	
		// configure response Bad Request
		fmt.Fprintf(w, responseNodeError, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest), r.Method + " " + r.URL.String())
				
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusBadRequest) + " " + http.StatusText(http.StatusBadRequest))
		fmt.Print(">>> ")
	}
}

//  resource: /X.X/controls
func HandlerControls(w http.ResponseWriter, r *http.Request)(){
	
	if(r.Method == get){
		// configure response OK
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		
		GPIOUpdate()
		json.NewEncoder(w).Encode(Io)
	
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK)  + " " + http.StatusText(http.StatusOK))
		fmt.Print(">>> ")
		
	} else if(r.Method == put){
		// changed I/O states
		var IoBuffer controls	// create a slice of pin
		json.NewDecoder(r.Body).Decode(&IoBuffer)
		
		// Configure the node with the new value
		GPIOConfig(IoBuffer)
		
		// configure response OK
		w.WriteHeader(http.StatusOK)
		
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK) + " " + http.StatusText(http.StatusOK))
		fmt.Print(">>> ")
		
	} else {	
		// configure response Bad Request
		fmt.Fprintf(w, responseNodeError, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest), r.Method + " " + r.URL.String())
					
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusBadRequest) + " " + http.StatusText(http.StatusBadRequest))
		fmt.Print(">>> ")
	}
}

//  resource: /X.X/configurations
func HandlerConfigurations(w http.ResponseWriter, r *http.Request)(){
		
	if(r.Method == get){
		// configure response OK
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(n.Config)
		
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK) + " " + http.StatusText(http.StatusOK))
		fmt.Print(">>> ")
		
	} else if(r.Method == put){
		//changed  the configuration
		var c configuration
		json.NewDecoder(r.Body).Decode(&c)
		
		// Configure the node with the new value
		ConfigureNode(c)
			
		// configure response OK
		w.WriteHeader(http.StatusOK)
		
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK) + " " + http.StatusText(http.StatusOK))
		fmt.Print(">>> ")
		
	} else {	
		// configure response Bad Request
		fmt.Fprintf(w, responseNodeError, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest), r.Method + " " + r.URL.String())
					
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusBadRequest) + " " + http.StatusText(http.StatusBadRequest))
		fmt.Print(">>> ")
	}
}

//  resrouce: /X.X/ping
func HandlerPing(w http.ResponseWriter, r *http.Request)(){		
	w.WriteHeader(http.StatusOK)
	
	log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusOK) + " " + http.StatusText(http.StatusOK))
	fmt.Print(">>> ")
}

//--	INITIALISATION		--//
// Define mux that contain map[resource]handlerfunc
var mux map[string]func(w http.ResponseWriter, r *http.Request)
							
/*															*
* Function: initServer()									*
* Param:													*
* Return:													*
* Description: 	initialize map[resource]handlerFunc and 	*
*				start the server							*
*															*/
func ServerInit()(){
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux[resourceNodesMeasurements] 		= HandlerMeasurements
	mux[resourceNodesControls]		 	= HandlerControls
	mux[resourceNodesConfigurations] 	= HandlerConfigurations
	mux[resourceNodesPing]			 	= HandlerPing
	
	// Start server
	serverNode.ListenAndServe()
}

/*															*
* Interface: (*requestHandler) ServeHTTP()					*
* Param:													* 	
*	- w http.ResponseWriter -> Response for the client		*
*	- r *http.Request -> Request from the client			*
* Return:													*
* Description: 	Search if an handler correspond with the	*
*				resource									*
*															*/
type requestHandler struct {}
func (*requestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)(){
	
	log.Println("\r- Request " +  r.Method + " " + r.URL.String() + " from " + r.RemoteAddr)
	
	var url string
	if(strings.Contains(r.URL.String(), "/?")){
		url = r.URL.String()[:strings.Index(r.URL.String(), "/?")]
	}else{
		url = r.URL.String()
	}
	
	// if a handler corresponding whith the request
	// Search handler only with resource no with the parameters
	if h, ok := mux[url]; ok {
		h(w, r)	// call fonction contain into the map 'mux'
		return
	// else return NotFound error
	} else {
		// configure response not found
		fmt.Fprintf(w, responseNodeError, http.StatusNotFound,
					http.StatusText(http.StatusNotFound), r.URL.String())
					
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusNotFound) + " " + http.StatusText(http.StatusNotFound))
		fmt.Print(">>> ")
	}
}
