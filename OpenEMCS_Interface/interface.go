package main

import (
	"os"
	"log"
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"github.com/jroimartin/gocui"
)

const (

	nodesView 		= "NodesView"
	ADCView 		= "ADCView"
	DACView			= "DACView"	
	IOView			= "IOView"	
	OUTView			= "OUTView"
	inputView		= "InputView"	

	bgColor = gocui.ColorGreen
	fgColor = gocui.ColorBlack
	
	title_Width				= 55
	title_Height			= 13
)

var nodeList 		*gocui.View
var ADCList		 	*gocui.View
var DACList		 	*gocui.View
var IOList		 	*gocui.View
var OUTList		 	*gocui.View
var input		 	*gocui.View
var information 	*gocui.View

var currentNodeID int
var currentNode int = 0
var currentADC 	int = 0
var currentDAC 	int = 0
var currentIO 	int = 0
var currentOUT 	int = 0

var maxNode	int = -1
var maxADC 	int = -1
var maxDAC 	int = -1
var maxIO 	int = -1
var maxOUT 	int = -1

type Mode struct {
	enable 		bool
	offsetX		int
	offsetY		int
	jumpUp		int
	jumpDown	int
}

var modeNode			= Mode{false, 0  , 0 , 8, 8}
var modeADC				= Mode{false, 0  , 0 , 6, 6}
var modeDAC				= Mode{false, 0  , 0 , 6, 6}
var modeIO				= Mode{false, 0  , 0 , 6, 6}
var modeOUT				= Mode{false, 0  , 0 , 5, 5}

var modeEditADC			= Mode{false, 0  , 0 , 0, 1}
var modeEditDAC			= Mode{false, 0  , 0 , 1, 1}
var modeEditIO			= Mode{false, 0  , 0 , 1, 1}
var modeEditOUT			= Mode{false, 0  , 0 , 1, 1}
var modeEditName		= Mode{false, 8  , 0 , 1, 1}
var modeEditInterval	= Mode{false, 12 , 0 , 1, 1}
var modeEditDescription	= Mode{false, 15 , 0 , 1, 0}

var modeEditEnable		= Mode{false, 10 , 0 , 0, 1}
var modeEditReference	= Mode{false, 13 , 0 , 1, 1}
var modeEditMode		= Mode{false, 8  , 0 , 1, 1}
var modeEditValue		= Mode{false, 9  , 0 , 1, 1}

var currentMode *Mode

// Lock the cursor, can't move Up and Down
var lockMove bool = false

var informationRef    string = "Reference: maximal voltage measured"
var informationInt    string = "Interval: interval between each measurements in second"
var informationEnable string = "Arrow Up: true\nArrow Down: false"
var informationValue  string = "Arrow Up: 1\nArrow Down: 0"
var informationMode   string = "Arrow Up: in\nArrow Down: out"

var logo []string = []string {	"  ____                   ______ __  __  _____  _____ 	\n",
								" / __ \\                 |  ____|  \\/  |/ ____|/ ____|	\n",
								"| |  | |_ __   ___ _ __ | |__  | \\  / | |    | (___  	\n",
								"| |  | | '_ \\ / _ \\ '_ \\|  __| | |\\/| | |     \\___ \\ 	\n",
								"| |__| | |_) |  __/ | | | |____| |  | | |____ ____) |	\n",
								" \\____/| .__/ \\___|_| |_|______|_|  |_|\\_____|_____/ 	\n",
								"       | |                                           	\n",
								"       |_|      										\n",
								"=======================================================\n",
								"= By Sven Ritz, Joel Bodenmann, Nathan Loretan        =\n",
								"=======================================================\n",}
		
var logfile *os.File
var log_filePath string 		= "log.txt"
		
func main(){
	var err error
	
	// Open Logfile
	os.Remove(log_filePath)
	logfile, err = os.OpenFile(log_filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file, %s\n", err)
	}
	log.SetOutput(logfile)
	
	DatabaseInit()
	
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.InputEsc = true
	g.Cursor = true
	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func initKeybindings(g *gocui.Gui) error {

	g.Cursor = true
	
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	
	// Define the handler for each keystrockes
	g.SetKeybinding(nodesView		, gocui.KeyArrowDown	, gocui.ModNone, MoveCursorDown) 
	g.SetKeybinding(ADCView			, gocui.KeyArrowDown	, gocui.ModNone, MoveCursorDown)
	g.SetKeybinding(DACView			, gocui.KeyArrowDown	, gocui.ModNone, MoveCursorDown)
	g.SetKeybinding(IOView			, gocui.KeyArrowDown	, gocui.ModNone, MoveCursorDown)
	g.SetKeybinding(OUTView			, gocui.KeyArrowDown	, gocui.ModNone, MoveCursorDown)
	g.SetKeybinding(inputView		, gocui.KeyArrowDown	, gocui.ModNone, ChangeValueDown)
	
	g.SetKeybinding(nodesView		, gocui.KeyArrowUp		, gocui.ModNone, MoveCursorUp) 
	g.SetKeybinding(ADCView			, gocui.KeyArrowUp		, gocui.ModNone, MoveCursorUp) 
	g.SetKeybinding(DACView			, gocui.KeyArrowUp		, gocui.ModNone, MoveCursorUp) 
	g.SetKeybinding(IOView			, gocui.KeyArrowUp		, gocui.ModNone, MoveCursorUp) 
	g.SetKeybinding(OUTView			, gocui.KeyArrowUp		, gocui.ModNone, MoveCursorUp) 
	g.SetKeybinding(inputView		, gocui.KeyArrowUp		, gocui.ModNone, ChangeValueUp) 
	
	g.SetKeybinding(nodesView		, gocui.KeyEnter		, gocui.ModNone, Edit) 
	g.SetKeybinding(ADCView			, gocui.KeyEnter		, gocui.ModNone, Edit) 
	g.SetKeybinding(DACView			, gocui.KeyEnter		, gocui.ModNone, Edit) 
	g.SetKeybinding(IOView			, gocui.KeyEnter		, gocui.ModNone, Edit) 
	g.SetKeybinding(OUTView			, gocui.KeyEnter		, gocui.ModNone, Edit) 
	g.SetKeybinding(inputView		, gocui.KeyEnter		, gocui.ModNone, Edit) 
	
	g.SetKeybinding(nodesView		, gocui.KeyEsc			, gocui.ModNone, Esc) 
	g.SetKeybinding(ADCView			, gocui.KeyEsc			, gocui.ModNone, Esc) 
	g.SetKeybinding(DACView			, gocui.KeyEsc			, gocui.ModNone, Esc) 
	g.SetKeybinding(IOView			, gocui.KeyEsc			, gocui.ModNone, Esc) 
	g.SetKeybinding(OUTView			, gocui.KeyEsc			, gocui.ModNone, Esc) 
	
	return nil
}
	
func ResetCursor(g *gocui.Gui) {
	g.SetCurrentView(nodesView)
	nodeList.SetCursor(0, 0)
	ADCList.SetCursor(0, 0)
	DACList.SetCursor(0, 0)
	IOList.SetCursor(0, 0)
	OUTList.SetCursor(0, 0)
	
	modeNode.Enable()
	
	if(input != nil) {
		g.DeleteView(inputView)
	}
	
	currentNode = 0
	currentADC = 0
	currentDAC = 0
	currentIO = 0
	currentOUT = 0
	
	lockMove = false
}
	
/*															*
* Function/Interface: layout
* Param:						
* Return:													
* Description: 	define and position the different view, this
* 				function is automatically called for refresh
* 				the terminal
*															*/
func layout(g *gocui.Gui) error {
	
	width, height := g.Size()

	if v, err := g.SetView("title", 0, 0, title_Width, title_Height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		
		buffer := bytes.NewBufferString(strings.Join(logo, ""))
		v.Write(buffer.Bytes())
	}

	if v, err := g.SetView("information", title_Width + 1, 0, width - 1, title_Height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		information = v
		information.Wrap 	= true
		information.Frame 			= false
	}	
	
	if v, err := g.SetView(nodesView, 0, title_Height + 1, width/3, height - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		nodeList = v
		nodeList.Highlight 			= true
		nodeList.Wrap	 			= false
		nodeList.SelBgColor 		= bgColor
		nodeList.SelFgColor 		= fgColor
		nodeList.Title 				= "Nodes list"
		nodeList.Frame 				= true
	
		g.SetCurrentView(nodesView)
		modeNode.Enable()
	}	
	
	if v, err := g.SetView(ADCView, width/3 + 1, title_Height + 1, 2*width/3, height - (height - title_Height)/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ADCList = v
		ADCList.Highlight 		= true
		nodeList.Wrap	 		= false
		ADCList.SelBgColor 		= bgColor
		ADCList.SelFgColor 		= fgColor
		ADCList.Title 			= "A/D Converter"
		ADCList.Frame 			= true
	}
	
	if v, err := g.SetView(DACView, width/3 + 1, height - (height - title_Height)/2 + 1, 2*width/3, height - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		DACList = v
		DACList.Highlight 		= true
		nodeList.Wrap	 		= false
		DACList.SelBgColor 		= bgColor
		DACList.SelFgColor 		= fgColor
		DACList.Title 			= "D/A Converter"
		DACList.Frame 			= true
	}
	
	if v, err := g.SetView(IOView, 2*width/3 + 1, title_Height + 1, width - 1, height - (height - title_Height)/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		IOList = v
		IOList.Highlight 		= true
		nodeList.Wrap	 		= false
		IOList.SelBgColor 		= bgColor
		IOList.SelFgColor 		= fgColor
		IOList.Title 			= "I/O"
		IOList.Frame 			= true
	}
	
	if v, err := g.SetView(OUTView, 2*width/3 + 1, height - (height - title_Height)/2 + 1, width - 1, height - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		OUTList = v
		OUTList.Highlight 		= true
		nodeList.Wrap	 		= false
		OUTList.SelBgColor 		= bgColor
		OUTList.SelFgColor 		= fgColor
		OUTList.Title 			= "OUTPUT"
		OUTList.Frame 			= true
	}

	DisplayNode(g)
		
	if(currentNodeID == -1) {
		ADCList.Clear()
		DACList.Clear()
		IOList.Clear()
		OUTList.Clear()
	} else {
		
		Io, err := GetNodeControls(currentNodeID)
		if(err != nil ){
			return nil
		}
		
		DisplayADC(currentNodeID, Io, g)
		DisplayDAC(currentNodeID, Io, g)
		DisplayIO(currentNodeID, Io, g)
		DisplayOUT(currentNodeID, Io, g)
	}
	
	return nil
}
	
/*															*
* Function/Interface: DisplayNode
* Param:						
* Return:													
* Description: 	Display the different nodes connected with
* 				the id, name, descrption and interval, and
* 				the link for ADC, DAC, IO, OUT
*															*/
func DisplayNode(g *gocui.Gui)(error) {
	var err error
	var list []string 
	
	nodes, err := DatabaseGetAllNodes()
	if(err != nil) {
		return err
	}
	
	maxNode = -1
	currentNodeID = -1
	
	for i := range nodes {
		
		// Get the ID of the selected node
		if(i == currentNode) {
			currentNodeID = nodes[currentNode].ID
		}
		
		list = append(list, fmt.Sprintf("Node ID: %d\n", nodes[i].ID))
		list = append(list, "  -> A/D converter\n")
		list = append(list, "  -> D/A converter\n")
		list = append(list, "  -> I/O \n")
		list = append(list, "  -> OUTPUT \n")
		list = append(list, fmt.Sprintf("  Name: %s\n", nodes[i].Config.Name))
		list = append(list, fmt.Sprintf("  Interval: %d\n", nodes[i].Config.Interval))
		list = append(list, fmt.Sprintf("  Description: %s\n", nodes[i].Config.Description))
		list = append(list, "\n")
		
		maxNode = i 
	}	
	
	if(maxNode == -1 || currentNodeID == -1) {
		ResetCursor(g)
	}
	
	// Print on the view
	buffer := bytes.NewBufferString(strings.Join(list, ""))
	nodeList.Clear()
	nodeList.Write(buffer.Bytes())
	
	return err
}

/*															*
* Function/Interface: DisplayADC
* Param:	
* 	id:	Id of the node that will configurate					
* Return:													
* Description: 	Display the ADC of the node with enable,
* 				description and reference
*															*/
func DisplayADC(id int, Io controls, g *gocui.Gui)(error) {
	var err error
	var list []string
		
	maxADC = -1
		
	for i := range Io.ADConvert {
		list = append(list, fmt.Sprintf("A/D Converter %d:\n", i))
		list = append(list, fmt.Sprintf("  Enable: %s\n", strconv.FormatBool(Io.ADConvert[i].Enable)))
		list = append(list, "  Value: " + strconv.FormatFloat(float64(Io.ADConvert[i].Value), 'f', -1, 32) + "\n")
		list = append(list, "  Reference: " + strconv.FormatFloat(float64(Io.ADConvert[i].Reference), 'f', -1, 32) + "\n")
		list = append(list, fmt.Sprintf("  Description: %s\n", Io.ADConvert[i].Description))
		list = append(list, "\n")
				
		maxADC = i 
	}
	
	if(maxNode == -1 || currentNodeID == -1) {
		ResetCursor(g)
	}
	
	// Print on the view
	buffer := bytes.NewBufferString(strings.Join(list, ""))
	ADCList.Clear()
	ADCList.Write(buffer.Bytes())
	
	return err
}

/*															*
* Function/Interface: DisplayDAC
* Param:		
* 	id:	Id of the node that will configurate					
* Return:													
* Description: 	Display the DAC of the node with enable, 
* 				description, reference and the value
*															*/
func DisplayDAC(id int, Io controls, g *gocui.Gui)(error) {
	var err error
	var list []string

	maxDAC = -1
	
	for i := range Io.DAConvert {
		list = append(list, fmt.Sprintf("D/A Converter %d:\n", i))
		list = append(list, fmt.Sprintf("  Enable: %s\n", strconv.FormatBool(Io.DAConvert[i].Enable)))		
		list = append(list, "  Value: " + strconv.FormatFloat(float64(Io.DAConvert[i].Value), 'f', -1, 32) + "\n")
		list = append(list, "  Reference: " + strconv.FormatFloat(float64(Io.DAConvert[i].Reference), 'f', -1, 32) + "\n")
		list = append(list, fmt.Sprintf("  Description: %s\n", Io.DAConvert[i].Description)) 
		list = append(list, "\n")
		
		maxDAC = i 		
	}
	
	if(maxNode == -1 || currentNodeID == -1) {
		ResetCursor(g)
	}
	
	// Print on the view
	buffer := bytes.NewBufferString(strings.Join(list, ""))
	DACList.Clear()
	DACList.Write(buffer.Bytes())
	
	return err
}

/*															*
* Function/Interface: DisplayIO
* Param:
* 	id:	Id of the node that will configurate							
* Return:													
* Description: 	Display the IO of the node with enable, 
* 				description, value(print only if mode = in)
* 				and mode
*															*/
func DisplayIO(id int, Io controls, g *gocui.Gui)(error) {
	var err error
	var list []string

	maxIO = -1
	
	for i := range Io.DigitalIO {
		list = append(list, fmt.Sprintf("I/O %d:\n", i))
		list = append(list, fmt.Sprintf("  Enable: %s\n", strconv.FormatBool(Io.DigitalIO[i].Enable)))
		
		// Don't print the value if mode is in
		if(Io.DigitalIO[i].Mode == "out") {
			list = append(list, fmt.Sprintf("  Value: %d\n", Io.DigitalIO[i].Value))
		} else {
			list = append(list, "  Value: \n")
		}
		
		list = append(list, fmt.Sprintf("  Mode: %s\n", Io.DigitalIO[i].Mode))
		list = append(list, fmt.Sprintf("  Description: %s\n", Io.DigitalIO[i].Description))
		list = append(list, "\n")
		
		maxIO = i 
	}
	
	if(maxNode == -1 || currentNodeID == -1) {
		ResetCursor(g)
	}
	
	// Print on the view
	buffer := bytes.NewBufferString(strings.Join(list, ""))
	IOList.Clear()
	IOList.Write(buffer.Bytes())
	
	return err
}

/*															*
* Function/Interface: DisplayOUT
* Param:
* 	id:	Id of the node that will configurate							
* Return:													
* Description: 	Display the outputs of the node with enable,
* 				descrption, value
*															*/
func DisplayOUT(id int, Io controls, g *gocui.Gui)(error) {
	var err error
	var list []string
		
	maxOUT = -1
		
	for i := range Io.DigitalOut {
		list = append(list, fmt.Sprintf("OUTPUT %d:\n", i))
		list = append(list, fmt.Sprintf("  Enable: %s\n", strconv.FormatBool(Io.DigitalOut[i].Enable)))
		list = append(list, fmt.Sprintf("  Value: %d\n", Io.DigitalOut[i].Value))
		list = append(list, fmt.Sprintf("  Description: %s\n", Io.DigitalOut[i].Description))			
		list = append(list, "\n")
		
		maxOUT = i 
	}
	
	if(maxNode == -1 || currentNodeID == -1) {
		ResetCursor(g)
	}
	
	// Print on the view
	buffer := bytes.NewBufferString(strings.Join(list, ""))
	OUTList.Clear()
	OUTList.Write(buffer.Bytes())
	
	return err
}

/*															*
* Function/Interface: ChangeValueDown
* Param:
* 	v:	the view that call the interrupt						
* Return:													
* Description: 	Modifiy the value for the information enable,
* 				value(only IO and OUT) and mode with 
* 				the arrow down
*															*/
func ChangeValueDown(g *gocui.Gui, v *gocui.View)(error) {
	
	// Modify the value with arrorDown for Information about value, mode and enable
	switch {
		case modeEditEnable.IsActive() :
			buffer := bytes.NewBufferString("false")
			input.Clear()
			input.Write(buffer.Bytes())
		
		case modeEditValue.IsActive() :
			buffer := bytes.NewBufferString("0")
			input.Clear()
			input.Write(buffer.Bytes())
		
		case modeEditMode.IsActive() :	
			buffer := bytes.NewBufferString("out")
			input.Clear()
			input.Write(buffer.Bytes())
	}
	
	return nil
}

/*															*
* Function/Interface: ChangeValueUp
* Param:
* 	v:	the view that call the interrupt						
* Return:													
* Description: 	Modifiy the value for the information enable,
* 				value(only IO and OUT) and mode with 
* 				the arrow up
*															*/
func ChangeValueUp(g *gocui.Gui, v *gocui.View)(error) {
	
	// Modify the value with arrorUp for Information about value, mode and enable
	switch {
		case modeEditEnable.IsActive() :
			buffer := bytes.NewBufferString("true")
			input.Clear()
			input.Write(buffer.Bytes())
		
		case modeEditValue.IsActive() :
			buffer := bytes.NewBufferString("1")
			input.Clear()
			input.Write(buffer.Bytes())
		
		case modeEditMode.IsActive() :	
			buffer := bytes.NewBufferString("in")
			input.Clear()
			input.Write(buffer.Bytes())
	}

	return nil
}

/*															*
* Function/Interface: MoveCursorDown
* Param:
* 	v:	the view that call the interrupt						
* Return:													
* Description: move the cursor and change the selected mode
*															*/
func MoveCursorDown(g *gocui.Gui, v *gocui.View)(error) {
	
	// Check if the cursor can be moved
	if  currentMode.jumpDown != 0 && !lockMove && maxNode != -1 && maxADC != -1 && maxDAC != -1 && maxIO != -1 && maxOUT != -1 { 

		rememberCurrentMode := currentMode
		
		// Move the cursor and change the mode
		if v.Name() == nodesView {
			switch {
				case modeNode.IsActive() :
					if currentNode == maxNode {
						return nil
					}
					currentNode++
					
				case modeEditADC.IsActive() :
					modeEditDAC.Enable()
				
				case modeEditDAC.IsActive() :
					modeEditIO.Enable()
				
				case modeEditIO.IsActive() :
					modeEditOUT.Enable()			
				
				case modeEditOUT.IsActive() :
					modeEditName.Enable()		
				
				case modeEditName.IsActive() :
					modeEditInterval.Enable()
				
				case modeEditInterval.IsActive() :
					modeEditDescription.Enable()
			
				default :
					return nil
			}
			
		} else if v.Name() == ADCView {
			switch {
				case modeADC.IsActive() :
					if currentADC == maxADC {
						return nil
					}
					currentADC++
					
				case modeEditEnable.IsActive() :
					modeEditValue.Enable()
					
				case modeEditValue.IsActive() :
					modeEditReference.Enable()
				
				case modeEditReference.IsActive() :
					modeEditDescription.Enable()
					
				default :
					return nil
			}
			
		} else if v.Name() == DACView {	
			switch {
				case modeDAC.IsActive() :
					if currentDAC == maxDAC {
						return nil
					}
					currentDAC++
				
				case modeEditEnable.IsActive() :
					modeEditValue.Enable()
				
				case modeEditValue.IsActive() :
					modeEditReference.Enable()
				
				case modeEditReference.IsActive() :
					modeEditDescription.Enable()
					
				default :
					return nil
			}
			
		} else if v.Name() == IOView {
			switch {
				case modeIO.IsActive() :
					if currentIO == maxIO {
						return nil
					}
					currentIO++
				
				case modeEditEnable.IsActive() :
					modeEditValue.Enable()
				
				case modeEditValue.IsActive() :
					modeEditMode.Enable()
				
				case modeEditMode.IsActive() :
					modeEditDescription.Enable()
					
				default :
					return nil
			}
			
		} else if v.Name() == OUTView {	
			switch {
				case modeOUT.IsActive() :
					if currentOUT == maxOUT {
						return nil
					}
					currentOUT++
	
				case modeEditEnable.IsActive() :
					modeEditValue.Enable()
				
				case modeEditValue.IsActive() :
					modeEditDescription.Enable()
					
				default :
					return nil
			}
		}
		
		_, y := v.Cursor()
		_, oy := v.Origin()
		
		if err := v.SetCursor(0, y + rememberCurrentMode.jumpDown); err != nil {
			v.SetOrigin(0, oy + rememberCurrentMode.jumpDown)
		}
	}

	return nil
}

/*															*
* Function/Interface: MoveCursorUp
* Param:
* 	v:	the view that call the interrupt						
* Return:													
* Description: move the cursor and change the selected mode
*															*/
func MoveCursorUp(g *gocui.Gui, v *gocui.View)(error) {

	// Check if the cursor can be moved
	if  currentMode.jumpUp != 0 && !lockMove && maxNode != -1 && maxADC != -1 && maxDAC != -1 && maxIO != -1 && maxOUT != -1 { 

		rememberCurrentMode := currentMode
		
		// Move the cursor and change the mode
		if v.Name() == nodesView {
			switch {
				case modeNode.IsActive() :
					currentNode--
					if currentNode < 0 {
						currentNode = 0
					}
				
				case modeEditDAC.IsActive() :
					modeEditADC.Enable()
				
				case modeEditIO.IsActive() :
					modeEditDAC.Enable()
				
				case modeEditOUT.IsActive() :
					modeEditIO.Enable()
				
				case modeEditName.IsActive() :
					modeEditOUT.Enable()
				
				case modeEditInterval.IsActive() :
					modeEditName.Enable()
								
				case modeEditDescription.IsActive() :
					modeEditInterval.Enable()
					
				default :
					return nil
					
			}	
				
		} else if v.Name() == ADCView {
			switch {
				case modeADC.IsActive() :
					currentADC--
					if currentADC < 0 {
						currentADC = 0
					}
				
				case modeEditValue.IsActive() :
					modeEditEnable.Enable()
				
				case modeEditReference.IsActive() :
					modeEditValue.Enable()
				
				case modeEditDescription.IsActive() :
					modeEditReference.Enable()
					
				default :
					return nil
			}
			
		} else if v.Name() == DACView {	
			switch {
				case modeDAC.IsActive() :
					currentDAC--
					if currentDAC < 0 {
						currentDAC = 0
					}
					
				case modeEditValue.IsActive() :
					modeEditEnable.Enable()
				
				case modeEditReference.IsActive() :
					modeEditValue.Enable()
				
				case modeEditDescription.IsActive() :
					modeEditReference.Enable()
					
				default :
					return nil
			}
			
		} else if v.Name() == IOView {
			switch {
				case modeIO.IsActive() :
					currentIO--
					if currentIO < 0 {
						currentIO = 0
					}
					
				case modeEditValue.IsActive() :
					modeEditEnable.Enable()
				
				case modeEditMode.IsActive() :
					modeEditValue.Enable()
				
				case modeEditDescription.IsActive() :
					modeEditMode.Enable()
					
				default :
					return nil
			}
			
		} else if v.Name() == OUTView {	
			switch {
				case modeOUT.IsActive() :
					currentOUT--
					if currentOUT < 0 {
						currentOUT = 0
					}
					
				case modeEditValue.IsActive() :
					modeEditEnable.Enable()
				
				case modeEditDescription.IsActive() :
					modeEditValue.Enable()
					
				default :
					return nil
			}
		}
		
		_, y := v.Cursor()
		_, oy := v.Origin()
		
		if err := v.SetCursor(0, y - rememberCurrentMode.jumpUp); err != nil && oy > 0 {
			v.SetOrigin(0, oy - rememberCurrentMode.jumpUp)
		}
	}

	return nil
}

/*															*
* Function/Interface: Esc
* Param:
* 	v:	the view that call the interrupt						
* Return:													
* Description: go out from a mode or view
*															*/
func Esc(g *gocui.Gui, v *gocui.View)(error) {
	var offset int
	var selectView	*gocui.View
	
	lockMove 	= false
	
	// Define the offset for position the cursor depending the View and the mode
	if v.Name() == nodesView {	
		
		selectView = nodeList	
		
		switch {	
			case modeEditADC.IsActive() :	
				modeNode.Enable()
				offset = 1		
			
			case modeEditDAC.IsActive() :
				modeNode.Enable()
				offset = 2		
			
			case modeEditIO.IsActive() :
				modeNode.Enable()
				offset = 3		
			
			case modeEditOUT.IsActive() :
				modeNode.Enable()
				offset = 4		
			
			case modeEditName.IsActive() :
				modeNode.Enable()
				offset = 5  
				 
			case modeEditInterval.IsActive() :
				modeNode.Enable()
				offset = 6
					 
			case modeEditDescription.IsActive() :
				modeNode.Enable()
				offset = 7
		}
	} else if v.Name() == ADCView {		
		
		selectView = ADCList		
		
		switch {
			case modeEditEnable.IsActive() :
				modeADC.Enable()
				offset = 1
				
			case modeEditValue.IsActive() :
				modeADC.Enable()
				offset = 2
				
			case modeEditReference.IsActive() :
				modeADC.Enable()
				offset = 3
				
			case modeEditDescription.IsActive() :
				modeADC.Enable()
				offset = 4
				
			case modeADC.IsActive() :
				// Reset the cursor
				currentADC = 0
				ADCList.SetCursor(0, 0)
				
				modeEditADC.Enable()
				g.SetCurrentView(nodesView)	
		}
	} else if v.Name() == DACView {		
		
		selectView = DACList
				
		switch {
			case modeEditEnable.IsActive() :
				modeDAC.Enable()
				offset = 1
				
			case modeEditValue.IsActive() :
				modeDAC.Enable()
				offset = 2
				
			case modeEditReference.IsActive() :
				modeDAC.Enable()
				offset = 3
				
			case modeEditDescription.IsActive() :
				modeDAC.Enable()
				offset = 4
				
			case modeDAC.IsActive() :
				// Reset the cursor
				currentDAC = 0
				DACList.SetCursor(0, 0)
			
				modeEditDAC.Enable()
				g.SetCurrentView(nodesView)	
		}
	} else if v.Name() == IOView {		
		
		selectView = IOList
		
		switch {	
			case modeEditEnable.IsActive() :
				modeIO.Enable()
				offset = 1
				
			case modeEditValue.IsActive() :
				modeIO.Enable()
				offset = 2
				
			case modeEditMode.IsActive() :
				modeIO.Enable()
				offset = 3
				
			case modeEditDescription.IsActive() :
				modeIO.Enable()
				offset = 4
				
			case modeIO.IsActive() :
				// Reset the cursor
				currentIO = 0
				IOList.SetCursor(0, 0)
			
				modeEditIO.Enable()
				g.SetCurrentView(nodesView)			
		}
	} else if v.Name() == OUTView {		
		
		selectView = OUTList
			
		switch {
			case modeEditEnable.IsActive() :
				modeOUT.Enable()
				offset = 1
				
			case modeEditValue.IsActive() :
				modeOUT.Enable()
				offset = 2
				
			case modeEditDescription.IsActive() :
				modeOUT.Enable()
				offset = 3
				
			case modeOUT.IsActive() :
				// Reset the cursor
				currentOUT = 0
				OUTList.SetCursor(0, 0)
			
				modeEditOUT.Enable()
				g.SetCurrentView(nodesView)	
		}
	}
	
	if selectView != nil {
		_, y := selectView.Cursor()	
		selectView.SetCursor(0, y - offset) 
	}
	
	return nil
}

/*															*
* Function/Interface: Edit
* Param:
* 	v:	the view that call the interrupt						
* Return:													
* Description: naviguate into the tree, selected the information
* 				that will be configurate and stat/stop editing
*															*/
func Edit(g *gocui.Gui, v *gocui.View)(error) {
	
	// The behavior changed depended the mode
	switch {
		case modeNode.IsActive() :
			// Can select the Node if no nodes are displayed
			if(maxNode == -1) {
				return nil
			}
			
			_, y := v.Cursor()
			v.SetCursor(0, y + 1)
			
			modeEditADC.Enable()
		
		case modeADC.IsActive() :
			fallthrough
			
		case modeDAC.IsActive() :
			fallthrough
				
		case modeIO.IsActive() :
			fallthrough
			
		case modeOUT.IsActive() :	
			// Can select the Node if no nodes are displayed		
			if(maxADC == -1 || maxDAC == -1 || maxIO == -1 || maxOUT == -1) {
				return nil
			}

			_, y := v.Cursor()
			v.SetCursor(0, y + 1)
			modeEditEnable.Enable()
		
		case modeEditADC.IsActive() :
			modeADC.Enable()
			g.SetCurrentView(ADCView)
		
		case modeEditDAC.IsActive() :
			modeDAC.Enable()
			g.SetCurrentView(DACView)
		
		case modeEditIO.IsActive() :
			modeIO.Enable()
			g.SetCurrentView(IOView)
		
		case modeEditOUT.IsActive() :
			modeOUT.Enable()
			g.SetCurrentView(OUTView)
				
		case modeEditName.IsActive() :
			newName, _ := EditMode(g, v, modeEditName.offsetX)
			
			// Clean the string
			newName = stripchars(newName, "\t\n\r\x00")
			newName = strings.Replace(newName, "\u0000", "", -1)
			
			if newName != "" { 
				n, _ := DatabaseGetNode(currentNodeID)
				n.Config.Name = newName
				PutNodeConfigurations(currentNodeID, n.Config)				
			}
				
		case modeEditDescription.IsActive() :
			newDescription, _ := EditMode(g, v, modeEditDescription.offsetX)
			
			// Clean the string
			newDescription = stripchars(newDescription, "\t\n\r\x00")
			newDescription = strings.Replace(newDescription, "\u0000", "", -1)
	
			if newDescription != "" {
				SetDescription(g, newDescription)
			}
	
		case modeEditInterval.IsActive() :
			// Print on the view
			information.Clear()
			information.Write([]byte(informationInt))		
		
			temp, _ := EditMode(g, v, modeEditInterval.offsetX)
			temp = stripchars(temp, "\t\n\r\x00 abcdefghijclmnopqrstuvwxyz")
			
			newInterval, err := strconv.Atoi(temp)		
			if err == nil {
				n, _ := DatabaseGetNode(currentNodeID)
				n.Config.Interval = newInterval
				PutNodeConfigurations(currentNodeID, n.Config)				
				information.Clear()
			}
		
		case modeEditEnable.IsActive() :
			// Print on the view
			information.Clear()
			information.Write([]byte(informationEnable))
					
			temp, _ := EditMode(g, v, modeEditEnable.offsetX)
			
			newEnable, err := strconv.ParseBool(temp)
			if err == nil {
				SetEnable(g, newEnable)
				information.Clear()
			}
			
		case modeEditReference.IsActive() :
			// Print on the view
			information.Clear()
			information.Write([]byte(informationRef))
			
			temp, _ := EditMode(g, v, modeEditReference.offsetX)
			temp = stripchars(temp, "\t\n\r\x00 abcdefghijclmnopqrstuvwxyz")
			
			newReference, err := strconv.ParseFloat(temp, 32) 
			if err == nil {
				SetReference(g, float32(newReference))
				information.Clear()
			}
	
		case modeEditMode.IsActive() :
			// Print on the view
			information.Clear()
			information.Write([]byte(informationMode))
		
			newMode, _ := EditMode(g, v, modeEditMode.offsetX)
		
			if strings.Compare(newMode, "in") == 0 || strings.Compare(newMode, "out") == 0 {
				SetMode(g, newMode)
				information.Clear()
			}
		
		case modeEditValue.IsActive() :
			// Print on the view
			information.Clear()
			information.Write([]byte(informationValue))
			
			if v.Name() == ADCView {
				return nil
			} else {
			
				temp, _ := EditMode(g, v, modeEditValue.offsetX)
				temp = stripchars(temp, "\t\n\r\x00 abcdefghijclmnopqrstuvwxyz")
				
				newValue, err := strconv.ParseFloat(temp, 32) 
				if err == nil {
					SetValue(g, float32(newValue))
					information.Clear()
				}
			}			
	}
	
	return nil
}	

/*															*
* Function/Interface: EditMode
* Param:
* 	v:	the view that call the interrupt						
* 	offsetX: offset for move the cursor
* Return:													
* Description: 	create a new view for editing the different
* 				information, the view is erase when it's finish
*															*/
var last string		// Save the name of the previous view
func EditMode(g *gocui.Gui, v *gocui.View, offsetX int)(string, error) {
	var err error
	var newData string
	
	// Exit editing
	if lockMove {	
		lockMove 	= false
		
		lastView, _ := g.View(last)
		lastView.Highlight = true
				
		_, cursorY := input.Cursor()
		input.SetCursor(1, cursorY)
		newData, _ = input.Line(cursorY)
		
		g.SetCurrentView(last)
		g.DeleteView(inputView)
		
	// Editing
	} else if !lockMove {
		lockMove 	= true		
		v.Highlight = false
		last = v.Name()
		
		x0, y0, x1, _, _ := g.ViewPosition(v.Name())
		_, cursorY := v.Cursor()
		
		// Create a View for editing
		if input, err = g.SetView(inputView, x0 + offsetX, y0 + cursorY, x1, y0 + cursorY + 2); err != nil {
			if err != gocui.ErrUnknownView {
				return newData, err
			}
			
			// Can write text for enable, mode, and binary value
			if !modeEditEnable.IsActive() && !modeEditMode.IsActive() && !(modeEditValue.IsActive() && (last == IOView || last == OUTView)) {
				input.Editable 	= true
				information.Clear()
			} 
			
			input.Frame = false
			
			// Write the current information on the input view
			str, _ := v.Line(cursorY)
			buffer := bytes.NewBufferString(str)
			input.Write(buffer.Bytes()[offsetX:])
			
			if _, err := g.SetCurrentView(inputView); err != nil {
				return newData, err
			}
		}
	}
	
	return newData, err
}

func SetDescription(g *gocui.Gui, description string)(){
	
	n, _ := DatabaseGetNode(currentNodeID)
	Io, err := GetNodeControls(currentNodeID)
	if err != nil {
		return
	}
	
	switch {
		case g.CurrentView() == nodeList :	
			n.Config.Description = description
			
		case g.CurrentView() == ADCList :
			for i := range Io.ADConvert {
				if currentADC == i {
					Io.ADConvert[i].Description = description
				}
			}
		
		case g.CurrentView() == DACList :
			for i := range Io.DAConvert {
				if currentDAC == i {
					Io.DAConvert[i].Description = description
				}
			}
		
		case g.CurrentView() == IOList :
			for i := range Io.DigitalIO {
				if currentIO == i {
					Io.DigitalIO[i].Description = description
				}
			}
		
		case g.CurrentView() == OUTList :
			for i := range Io.DigitalOut {
				if currentOUT == i {
					Io.DigitalOut[i].Description = description
				}
			}
	}
	
	// Send the new configuration to the node
	PutNodeConfigurations(currentNodeID, n.Config)		
	
	// Send the new controls to the node
	PutNodeControls(currentNodeID, Io)	
}

func SetEnable(g *gocui.Gui, enable bool)(){

	Io, err := GetNodeControls(currentNodeID)
	if err != nil {
		return
	}

	// Which device is selected
	switch {	
		case g.CurrentView() == ADCList :
			// Check which ADC must be set
			for i := range Io.ADConvert {
				if currentADC == i {
					Io.ADConvert[i].Enable = enable
				}
			}
		
		case g.CurrentView() == DACList :
			// Check which DAC must be set
			for i := range Io.DAConvert {
				if currentDAC == i {
					Io.DAConvert[i].Enable = enable
				}
			}
			
		case g.CurrentView() == IOList :
			// Check which IO must be set
			for i := range Io.DigitalIO {
				if currentIO == i {
					Io.DigitalIO[i].Enable = enable
				}
			}
			
		case g.CurrentView() == OUTList :
			// Check which Output must be set
			for i := range Io.DigitalOut {
				if currentOUT == i {
					Io.DigitalOut[i].Enable = enable
				}
			}	
	}

	// Send the new controls to the node
	PutNodeControls(currentNodeID, Io)	
}

func SetReference(g *gocui.Gui, reference float32)(){

	Io, err := GetNodeControls(currentNodeID)
	if err != nil {
		return
	}

	// Which device is selected
	switch {
		case g.CurrentView() == ADCList :
			// Check which ADC must be set
			for i := range Io.ADConvert {
				if currentADC == i {
					Io.ADConvert[i].Reference = reference
				}
			}
				
		case g.CurrentView() == DACList :
			// Check which DAC must be set
			for i := range Io.DAConvert {
				if currentDAC == i {
					Io.DAConvert[i].Reference = reference
				}
			}	
	}

	// Send the new controls to the node
	PutNodeControls(currentNodeID, Io)	
}

func SetValue(g *gocui.Gui, value float32)(){

	Io, err := GetNodeControls(currentNodeID)
	if err != nil {
		return
	}

	// Which device is selected
	switch {
		case g.CurrentView() == DACList :
			// Check which DAC must be set
			for i := range Io.DAConvert {
				if currentDAC == i {
					Io.DAConvert[i].Value = value
				}
			}
		
		case g.CurrentView() == IOList :
			// Check which IO must be set
			for i := range Io.DigitalIO {
				if currentIO == i && Io.DigitalIO[i].Mode == out {
					
					// only binary value
					if value >= 1 {
						value = 1
					} else {
						value = 0
					}
					Io.DigitalIO[i].Value = int(value)
				}
			}
			
		case g.CurrentView() == OUTList :
			// Check which Output must be set
			for i := range Io.DigitalOut {
				if currentOUT == i {
						
					// only binary value
					if value >= 1 {
						value = 1
					} else {
						value = 0
					}
					
					Io.DigitalOut[i].Value = int(value)
				}
			}
	}
	
	// Send the new controls to the node
	PutNodeControls(currentNodeID, Io)	
}

func SetMode(g *gocui.Gui, mode string)(){

	Io, err := GetNodeControls(currentNodeID)
	if err != nil {
		return
	}

	// Check which IO must be set
	for i := range Io.DigitalIO {
		if currentIO == i {
			Io.DigitalIO[i].Mode = mode
		}
	}
	
	PutNodeControls(currentNodeID, Io)	
}

func(m *Mode) Enable() {
	
	m.enable = true
	
	// Disable the old Mode
	if(currentMode != nil && currentMode != m) {
		currentMode.enable = false
	}
	
	currentMode = m
}

func(m *Mode) IsActive()(bool) {
	return m.enable

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

