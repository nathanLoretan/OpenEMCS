package main

import (
	"fmt"
	"bytes"
	"time"
	"errors"
	"net/http"
	"encoding/json"
)

var clientNode = &http.Client{
	Timeout: time.Second * 10,
}
	
/*															
* Function/Interface: 	GetNodeMeasurements					
* Param:													
	- id: ID of the nodes
	- nbrMeasurements: number of measurement that the node must return											 
* Return:													
	- Measurements: a structur with all measurements
	- error: the error that occures with the request
* Description:												
*															*/
func GetNodeMeasurements(id int, nbrMeasurements int)([]measurement, error){
	var measurements []measurement = nil
	var err error
	
	ip, err := DatabaseGetNodeIP(id)
	if ip != "" && err == nil {
		
		// Send request and get response
		req, _ := http.NewRequest(get, resourceHttp + ip + nodePort + resourceNodesMeasurements + fmt.Sprintf(paramMeasurements, nbrMeasurements), nil)
		resp, err := clientNode.Do(req)
		if err != nil{
			return measurements, err
		}
		// Close network connection
		defer resp.Body.Close()		
		
		// If no error and response OK
		if(err == nil && resp.StatusCode == http.StatusOK){			
			json.NewDecoder(resp.Body).Decode(&measurements)
		
		// If error from http request
		}else if(err == nil && resp.StatusCode != http.StatusOK){	
			// create an error for status code
			err = errors.New("Error: " + fmt.Sprintf("%d", resp.StatusCode) + " " + http.StatusText(resp.StatusCode))
		}
	}
		
	return  measurements, err
}

/*															
* Function/Interface: 	GetNodeNbrMeasurements					
* Param:													
	- id: ID of the nodes
* Return:													
	- int: the number of measurements into the node
	- error: the error that occures with the request
* Description:												
*															*/
func GetNodeNbrMeasurements(id int)(int, error){
	var nbrMeasurements int
	var err error
	
	ip, err := DatabaseGetNodeIP(id)
	if ip != "" && err == nil {
		
		// Send request and get response
		req, _ := http.NewRequest(get, resourceHttp + ip + nodePort + resourceNodesMeasurements + paramAll, nil)
		resp, err := clientNode.Do(req)
		if err != nil{
			return nbrMeasurements, err
		}
		// Close network connection
		defer resp.Body.Close()		
		
		// If no error and response OK
		if(err == nil && resp.StatusCode == http.StatusOK){			
			json.NewDecoder(resp.Body).Decode(&nbrMeasurements)
		
		// If error from http request
		}else if(err == nil && resp.StatusCode != http.StatusOK){	
			// create an error for status code
			err = errors.New("Error: " + fmt.Sprintf("%d", resp.StatusCode) + " " + http.StatusText(resp.StatusCode))
		}
	}
		
	return  nbrMeasurements, err
}

/*															
* Function/Interface: 		DeleteNodeMeasurements			
* Param:													
	- id: ID of the nodes 
* Return:													
	- error: the error that occures with the request
* Description: 												
*															*/
func DeleteNodeMeasurements(id int)(error){
	var err error
	
	ip, err := DatabaseGetNodeIP(id)
	if ip != "" && err == nil {
		
		// Send request and get response
		req, _ := http.NewRequest(delete, resourceHttp + ip + nodePort + resourceNodesMeasurements, nil)
		resp, err := clientNode.Do(req)
		if err == nil {
			return err
		}
		// Close network connection
		defer resp.Body.Close()		
		
		// If error from http request
		if(err == nil && resp.StatusCode != http.StatusOK){		
			// create an error for status code
			err = errors.New("Error: " + fmt.Sprintf("%d", resp.StatusCode) + " " + http.StatusText(resp.StatusCode))
		}
	}
	
	return err
}

/*															
* Function/Interface: 	GetNodeControls						
* Param:													
	- id: ID of the nodes 
* Return:													
	- Control: a structur with current control's configurations of the node
	- error: the error that occures with the request
* Description: 												
*															*/
func GetNodeControls(id int)(controls, error){
	var Io controls
	var err error	
	
	ip, err := DatabaseGetNodeIP(id)
	if ip != "" && err == nil {
		
		// Send request and get response
		req, _ := http.NewRequest(get, resourceHttp + ip + nodePort + resourceNodesControls, nil)
		resp, err := clientNode.Do(req)
		if err != nil {
			return Io, err
		}
		// Close network connection
		defer resp.Body.Close()		
		
		// If no error and response OK
		if(err == nil && resp.StatusCode == http.StatusOK){			
			json.NewDecoder(resp.Body).Decode(&Io)
		
		// If error from http request
		} else if(err == nil && resp.StatusCode != http.StatusOK){	
			// create an error for status code
			err = errors.New("Error: " + fmt.Sprintf("%d", resp.StatusCode) + " " + http.StatusText(resp.StatusCode))
		}
	}
	
	return Io, err
}

/*															
* Function/Interface: 		PutNodeControls					
* Param:													
	- id: ID of the nodes 
	- Control: a structur with new control's configurations
* Return:													
	- error: the error that occures with the request
* Description: 												
*															*/
func PutNodeControls(id int, Io controls)(error){
	var err error
	
	// Encode the body for request
	body := new(bytes.Buffer)
    json.NewEncoder(body).Encode(Io)

	ip, err := DatabaseGetNodeIP(id)
	if ip != "" && err == nil {
		
		// Send request and get response
		req, _ := http.NewRequest(put, resourceHttp + ip + nodePort + resourceNodesControls, body)
		resp, err := clientNode.Do(req)
		if err != nil{
			return err
		}
		// Close network connection
		defer resp.Body.Close()		
		
		// If error from http request
		if(err == nil && resp.StatusCode != http.StatusOK){		
			// create an error for status code
			err = errors.New("Error: " + fmt.Sprintf("%d", resp.StatusCode) + " " + http.StatusText(resp.StatusCode))
		}
	}
	
	return err
}

/*															
* Function/Interface: 		GetNodeConfigurations			
* Param:													
	- id: ID of the nodes 
* Return:													
	- Configurations: structur with the current node's configurations
	- error: the error that occures with the request
* Description: 												
*															*/
func GetNodeConfigurations(id int)(configuration, error){
	var configNode configuration
	var err error
	
	ip, err := DatabaseGetNodeIP(id)
	if ip != "" && err == nil {

		// Send request and get response
		req, _ := http.NewRequest(get, resourceHttp + ip + nodePort + resourceNodesConfigurations, nil)
		resp, err := clientNode.Do(req)
		if err != nil {
			return configNode, err
		}
		// Close network connection
		defer resp.Body.Close()		
		
		// If no error and response OK
		if(err == nil && resp.StatusCode == http.StatusOK){			
			json.NewDecoder(resp.Body).Decode(&configNode)
			
		// If error from http request
		}else if(err == nil && resp.StatusCode != http.StatusOK){	
			// create an error for status code
			err = errors.New("Error: " + fmt.Sprintf("%d", resp.StatusCode) + " " + http.StatusText(resp.StatusCode))
		}
	}
	
	return configNode, err
}

/*															
* Function/Interface: 	PutNodeConfigurations				
* Param:													
	- Configurations: structur with the new node's configurations
	- id: ID of the nodes 
* Return:													
	- error: the error that occures with the request
* Description: 												
*															*/
func PutNodeConfigurations(id int, c configuration)(error){
	var err error
	
	// Encode the body for request
	body := new(bytes.Buffer)
    json.NewEncoder(body).Encode(&c)
	
	ip, err := DatabaseGetNodeIP(id)
	if ip != "" && err == nil {
		
		// Send request and get response
		req, _ := http.NewRequest(put, resourceHttp + ip + nodePort + resourceNodesConfigurations, body)
		resp, err := clientNode.Do(req)
		if err != nil{
			return err
		}
		// Close network connection
		defer resp.Body.Close()		
		
		// If error from http request
		if(err == nil && resp.StatusCode != http.StatusOK){		
			// create an error for status code
			err = errors.New("Error: " + fmt.Sprintf("%d", resp.StatusCode) + " " + http.StatusText(resp.StatusCode))
		} else {
			_, err = DatabaseSaveNode(ip, node{id, c})	// save the new configuration acknowledged by the node
		}
	}
	
	return err
}
