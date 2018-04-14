package main

var serverIp string		= "127.0.0.1"

const(
	version 		= "1.0"
	nodePort		= ":8001"
	serverPort 		= ":8000"

	// Verbex
	get 	= "GET"
	put 	= "PUT"
	post	= "POST"
	delete 	= "DELETE"
	
	// Define the resource
	resourceHttp 	= "http://"
	resourceNodesMeasurements 	= "/" + version + "/measurements"
	resourceNodesControls 		= "/" + version + "/controls"
	resourceNodesConfigurations	= "/" + version + "/configurations"
	resourceNodesPing 			= "/" + version + "/ping"	
	resourcesNodes				= "/" + version + "/nodes"	
	
	// Define the response server->nodes/users
	responseServerError		= "Server Error: %v %v resource: %v"
	responseNodeError		= "Node Error: %v %v resource: %v"
	
	// Define parameters
	paramMeasurements		= "/?nbrMeasurements=%d"
	paramAll				= "/?all"
)

// Constant for the Pins
const(
	nbrAdc 	= 2
	nbrDac 	= 2
	nbrOut 	= 2
	nbrIO 	= 2
	
	out = "out"
	in 	= "in"	
)

type converter struct {
	Enable			bool		`json:"Enable"`
	Description		string		`json:"Description"`
	Value 			float32		`json:"Value"`
	Reference 		float32		`json:"Reference"`
}

type o struct {
	Enable			bool	`json:"Enable"`
	Description		string	`json:"Description"`
	Value 			int		`json:"Value"`
}

type io struct {
	Enable			bool	`json:"Enable"`
	Description		string	`json:"Description"`
	Mode			string	`json:"Mode"`
	Value 			int		`json:"Value"`
}

type controls struct {
	DigitalOut 	[nbrOut]o    		`json:"DigitalOutput"`
	DigitalIO 	[nbrIO]io 			`json:"DigitalIO"`
	ADConvert	[nbrAdc]converter 	`json:"ADConvert"`
	DAConvert	[nbrDac]converter 	`json:"DAConvert"`
}

type configuration struct {
	Name 		string	 `json:"NodeName"`
	Description string   `json:"Description"`
	Interval 	int	 	 `json:"Interval"`
}

type measurement struct {
	V			float32		`json:"V"`
	Vrms		float32		`json:"Vrms"`
	I			float32		`json:"I"`
	Irms		float32		`json:"Irms"`
	Time		int64 		`json:"Time"`
}

type node struct {
	ID			int 	`json:"ID"`
	Config		configuration
}
