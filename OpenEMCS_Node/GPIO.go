package main

/*
#cgo LDFLAGS: -lwiringPi
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <wiringPi.h>
#include <wiringPiSPI.h>

#define ADC	0
#define DAC	1

void print(char* s) {
    printf("%s", s);
}

void setup() {
 	printf("Setup\n");
 	wiringPiSetupGpio();
 	
 	// SPI channel 0, Clock 500kHz
 	// Spi must be enable with command: 	-> raspi-config
 	// Or enable dtparam=spi=on into the file /boot/config.txt
	wiringPiSPISetup (0, 100000) ;   			
	wiringPiSPISetup (1, 100000) ;   			
}

void writePin_ON(int pin) {
	digitalWrite(pin, 1);
}

void writePin_OFF(int pin) {
	digitalWrite(pin, 0);
}

void modePin_IN(int pin) {
	pinMode(pin, INPUT);
}

void modePin_OUT(int pin) {
	pinMode(pin, OUTPUT);
}

int readPin(int pin) {
	return digitalRead(pin);
}

int readAD(int channel) {
	
	if(channel > 7) {
		return 0;
	}
	
	// Send the MSB first, also send at first buffer bit 7.
	// MOSI 	Start | SGL/DIFF | D2 | D1 | D0 | 0 | 0 | 0  | 0  | ...		 | 0  | 0  | 0  | ...
	// MISO											| 0 | B9 | B8 | B7 | ... | B1 | B0 | B1 | B2 | ...
	// Buffer	|				Buffer[0]					 |		Buffer[1]	  |		Buffer[2]					
 	uint8_t buffer[3] = {0, 0, 0}; 
	buffer[0] += 1 << 7;					// Start bit	
	buffer[0] += 1 << 6;					// mode Single
	buffer[0] += (channel & 0x04) << 3;		// D2
	buffer[0] += (channel & 0x02) << 3;		// D1
	buffer[0] += (channel & 0x01) << 3;		// D0
	
	// Send first buffer[0] then buffer[1] then buffer[2]
	// Save first buffer[0] then buffer[1] then buffer[2]
	if(wiringPiSPIDataRW(ADC, buffer, 3) < 0) {
		return -1;
	}

	//		   |		B9			   |   |	B8...B1	  |   | 		B0			  |	
	int data = ((buffer[0] & 0x01) << 9) + (buffer[1] << 1) + ((buffer[2] & 0x80) >> 7); 

	return data;
}

void writeDA(int channel, int data) {
	
	if(channel > 1) {
		return;
	}
	
	// Send the MSB first, also send at first buffer bit 7.
	// MOSI 	DAa/b | BUF | /GA | /SHDN | D9 | D8 | D7 | D6 | ... | D1 | D0 | XX | XX |
	// Buffer	|				  Buffer[0]				      |			Buffer[1]	    |
	uint8_t buffer[2] = {0, 0}; 
	buffer[0] += (channel & 0x01) << 7;		// DA a or b	
	buffer[0] += 0 << 6;					// Unbufferd mode
	buffer[0] += 1 << 5;					// Vref
	buffer[0] += 1 << 4;					// Active mode, Vout is available 
	buffer[0] += (data & 0x200) >> 6;		// D9
	buffer[0] += (data & 0x100) >> 6;		// D8
	buffer[0] += (data & 0x080) >> 6;		// D7
	buffer[0] += (data & 0x040) >> 6;		// D6
	buffer[1] += (data & 0x020) << 2;		// D5
	buffer[1] += (data & 0x010) << 2;		// D4
	buffer[1] += (data & 0x008) << 2;		// D3
	buffer[1] += (data & 0x004) << 2;		// D2
	buffer[1] += (data & 0x002) << 2;		// D1
	buffer[1] += (data & 0x001) << 2;		// D0
	
	// Send first buffer[0] then buffer[1] 
	wiringPiSPIDataRW(DAC, buffer, 2);
} 
*/
import "C"

const(
	VChannel 	= 0
	VrmsChannel	= 1
	IChannel	= 2
	IrmsChannel	= 3
	
	convertResolution 		= 1024 
	measurementReference	= 230
) 

var	  OutPin  	= [nbrOut]int{15, 14}
var   IOPin 	= [nbrIO]int{23, 18}
var   ADChannel = [nbrAdc]int{4, 5}
var   DAChannel = [nbrDac]int{0, 1}


func GPIOSetup()(){
	C.setup();
}

func WritePin(pin int, value int)(){
	// Test if gpio is set mode out
	if(value == 0) {		
		C.writePin_OFF(C.int(pin) )
	} else {
		C.writePin_ON(C.int(pin) )
	}
}

func ModePin(pin int, mode string)(){
	if(mode == in) {
		C.modePin_IN(C.int(pin))
	} else {
		C.modePin_OUT(C.int(pin))
	}
}

func ReadPin(pin int)(int){
	var state C.int = 0
	state = C.readPin(C.int(pin))
	return int(state)
}

func ReadAD(channel int)(int){
	var data C.int = C.readAD(C.int(channel))
	return int(data)
}

func WriteDA(channel int, data int)(){
	C.writeDA(C.int(channel), C.int(data))
}

