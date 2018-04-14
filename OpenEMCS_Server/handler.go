// Server handler
package main

import (
	"log"
	"fmt"
	"time"
	"strings"
	"net/http"	
	"encoding/json"
)

// Define client and server for send/Receive request
var serverNode = http.Server{
	Addr: serverPort,
	Handler: &requestHandler{},
	ReadTimeout:  time.Second * 10,
	WriteTimeout: time.Second * 10,
}
						
//--	HANDLER		--//
//  resrouce: /X.X/nodes
func HandlerNode(w http.ResponseWriter, r *http.Request)(){
	
	if(r.Method == get){
		// Get the node frome the database
		nodes, err := DatabaseGetAllNodes()	
		
		if err == nil {
			// configure response OK
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nodes)
			
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
			
	} else if(r.Method == post){
		//changed  the delay between each measurements
		var n node	
		json.NewDecoder(r.Body).Decode(&n)
		id, err := DatabaseSaveNode(r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")], n)

		if err == nil {
			n.ID = id
			
			// configure response OK
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(n)	// return the new id of the node
		
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
		
	} else if(r.Method == delete){
		
		// Delete the node from the database
		err := DatabaseDeleteNode(r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")])
		
		if err == nil {
			// configure response ok
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
		// configure response Bad request
		fmt.Fprintf(w, responseServerError, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest), r.Method + " " + r.URL.String())
				
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusBadRequest) + " " + http.StatusText(http.StatusBadRequest))
		fmt.Print(">>> ")
	}
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
	mux[resourcesNodes] = HandlerNode
	
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
	
	// if a handler corresponfing whith the request
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)	// call fonction contain into the map 'mux'
		return
	// else return NotFound error
	} else {
		// configure response not found
		fmt.Fprintf(w, responseServerError, http.StatusNotFound,
					http.StatusText(http.StatusNotFound), r.URL.String())
					
		log.Println("\r\t Response " + fmt.Sprintf("%d", http.StatusNotFound) + " " + http.StatusText(http.StatusNotFound))
		fmt.Print(">>> ")
	}
}
