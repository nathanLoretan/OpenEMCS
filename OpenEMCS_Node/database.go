// Node database
package main

import(
	"log"
	"errors"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var measurementDB_filePath string 	= "../Database/measurementsDB.db"
var configurationDB_filePath string = "../Database/configurationDB.db"
var controlsDB_filePath string 		= "../Database/controlsDB.db"
var db_measurements *sql.DB
var db_configuration *sql.DB
var db_controls *sql.DB

/*															*
* Function/Interface: dataBaseInit
* Param:						
* Return:													
* Description: 	This function create the dataBases if any
*				others exist
*															*/
func DatabaseInit()(error){	
	var err error
	
	// Open a connection with the database
	db_measurements, err = sql.Open("sqlite3", measurementDB_filePath)
	db_configuration, err = sql.Open("sqlite3", configurationDB_filePath)
	db_controls, err = sql.Open("sqlite3", controlsDB_filePath)
	
	// close the database
	defer db_measurements.Close()
	defer db_configuration.Close()
	defer db_controls.Close()
	
	// Create the database if any database exists
	_, err = db_measurements.Exec(	"CREATE TABLE IF NOT EXISTS `measurements`(" +
									"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
									"`time` INTEGER NULL," + 
									"`v` REAL NULL," +
									"`vrms` REAL NULL," +
									"`i` REAL NULL," +
									"`irms` REAL NULL);")
									
	_, err = db_configuration.Exec(	"CREATE TABLE IF NOT EXISTS `configuration`(" +
									"`id` INTEGER NULL," +
									"`name`  VARCHAR(20) PRIMARY KEY NULL," +
									"`description` VARCHAR(100) NULL," +
									"`interval` REAL NULL);")
									
	_, err = db_controls.Exec(	"CREATE TABLE IF NOT EXISTS `controls`(" +
									"`pinCh` INTEGER NULL," +
									"`type` VARCHAR(3) NULL," + 
									"`enable` BOOLEAN NULL," +
									"`description` VARCHAR(100) NULL," +
									"`value` REAL NULL," +
									"`reference` REAL NULL," +
									"`mode` VARCHAR(3) NULL);")
									
	// close the database
	db_measurements.Close()
	db_configuration.Close()
	db_controls.Close()
	
	db_measurements, err = sql.Open("sqlite3", measurementDB_filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	// close the database
	defer db_measurements.Close()
	db_configuration, err = sql.Open("sqlite3", configurationDB_filePath)
	
	if err != nil {
		log.Println(err)
		return err
	}
	// close the database
	defer db_configuration.Close()
	
	db_controls, err = sql.Open("sqlite3", controlsDB_filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	// close the database
	defer db_controls.Close()
	
	return err
}

/*															*
* Function/Interface: databaseGetMeasurements
* Param:					
*	nbrMeasurements: 	the number of measurements that must be
*						taken from the database
* Return:		
*	measurement[]: slide with the mesaurements											
* Description: 	Get some measurements from the database
*															*/
func DatabaseGetMeasurements(nbrMeasurements int) ([]measurement, error){
	var err error
	var measurements []measurement

	n, err := DatabaseGetNbrMeasurements()
	if(err != nil || n == 0) {
		if(n == 0) {
			err = errors.New("Error: no measurements!!")
		}
		return measurements, err
	}

	// if the number of measurements within the database is lower
	// than the number queried	
	if(n < nbrMeasurements) {
		nbrMeasurements = n
	}
	
	// Open a connection with the database
	db_measurements, err = sql.Open("sqlite3", measurementDB_filePath)
	if err != nil {
		return measurements, err
	}
	// close the database
	defer db_measurements.Close()
	
	// Prepare a statement
	stmt, err := db_measurements.Prepare("SELECT * FROM measurements LIMIT ? OFFSET (SELECT COUNT(*) FROM measurements)-?;")
	if err != nil {
		return measurements, err
	}
	
	// Execute the statement
	rows, err := stmt.Query(nbrMeasurements, nbrMeasurements)
	if err != nil{
		return measurements, err
	}
	
	// Add all measurements to the slide measurements
	for rows.Next(){
		var id int
		var m measurement
		err = rows.Scan(&id, &m.Time, &m.V, &m.Vrms, &m.I, &m.Irms);
		if err != nil{
			return measurements, err
		}
		measurements = append(measurements, m)
	}

	return measurements, err
}

/*															*
* Function/Interface: databaseGetMeasurements
* Param:					
* Return:		
*	int: the number of measurements into the database											
* Description: 	Get the numbers of measurements contain within
				the database
*															*/
func DatabaseGetNbrMeasurements()(int, error) {
	var err error
	var nbrMeasurements int

	// Open a connection with the database
	db_measurements, err = sql.Open("sqlite3", measurementDB_filePath)
	if err != nil {
		return nbrMeasurements, err
	}
	// close the database
	defer db_configuration.Close()
	
	// Get ip address at the column indicate by id
	row := db_measurements.QueryRow("SELECT COUNT(*) FROM measurements;")
	err = row.Scan(&nbrMeasurements)
	
	return nbrMeasurements, err
}

/*															*
* Function/Interface: databaseDeleteMeasurements
* Param:						
* Return:													
* Description: 	Delete all measurements contained into the 
*				database
*															*/
func DatabaseDeleteMeasurements()(error){
	var err error

	// Open a connection with the database
	db_measurements, err = sql.Open("sqlite3", measurementDB_filePath)
	if err != nil {
		return err
	}
	// close the database
	defer db_measurements.Close()

	// Delete all value from database
	_, err = db_measurements.Exec("DELETE FROM measurements;")
	
	return err
}

/*															*
* Function/Interface: databaseAddMeasurements
* Param:						
*	m: measurement struct with time, voltage, current
* Return:													
* Description: 	Add one measurements to the database
*															*/
func DatabaseAddMeasurements(m measurement)(error){
	var err error

	// Open a connection with the database
	db_measurements, err = sql.Open("sqlite3", measurementDB_filePath)
	if err != nil {
		return err
	}
	// close the database
	defer db_measurements.Close()

	// Prepare a statement
	stmt, err := db_measurements.Prepare("INSERT INTO measurements(time, v, vrms, i, irms) values(?, ?, ?, ?, ?);")
	if err != nil {
		return err
	}
	
	// Execute the statement
	_, err = stmt.Exec(m.Time, m.V, m.Vrms, m.I, m.Irms)

	return err
}

/*															*
* Function/Interface: DatabaseCreateNode
* Param:						
* Return:						
*	node: the new node created in the database
* Description: 	Create a node with default value and add it
				in the database. The database must contain
				only one raw.
*															*/
func DatabaseCreateNode()(node, error){
	var err error
	var n node
	
	n.ID = 0
	n.Config.Name = "Unknown"
	n.Config.Description = "Unknown"
	n.Config.Interval = defaultInterval

	// Prepare a statement
	stmt, err := db_configuration.Prepare("INSERT INTO configuration(id, name, description, interval) VALUES(?, ?, ?, ?);")
	if err != nil {
		return n, err
	}
	
	// Execute the statement
	_, err = stmt.Exec(n.ID, n.Config.Name, n.Config.Description, n.Config.Interval)
	
	return n, err
		
}

/*															*
* Function/Interface: DatabaseGetNode
* Param:						
* Return:						
*	node: the node saved in the database
* Description: 	return the information of the node saved in 
				database. Call DatabaseCreateNode() if no
				one have already been created
*															*/
func DatabaseGetNode()(node, error){
	var err error
	var n node
	
	// Open a connection with the database
	db_configuration, err = sql.Open("sqlite3", configurationDB_filePath)
	if err != nil {
		return n, err
	}
	// close the database
	defer db_configuration.Close()
	
	// Get ip address at the column indicate by id
	row := db_configuration.QueryRow("SELECT id, name, description, interval FROM configuration;")
	err = row.Scan(&n.ID, &n.Config.Name, &n.Config.Description, &n.Config.Interval)
	
	// If it's the first connection of the ndoe
	if(err != nil){	
		n, err = DatabaseCreateNode();
		return n, err
	}
	
	return n, err
}

/*															*
* Function/Interface: DatabaseGetNode
* Param:		
*	n : node with the new informations
* Return:						
* Description:	change the informations about the node
*
*															*/
func DatabaseUpdateNode(n node)(error){
	var err error
	
	// Open a connection with the database
	db_configuration, err = sql.Open("sqlite3", configurationDB_filePath)
	if err != nil {
		return err
	}
	// close the database
	defer db_configuration.Close()

	// Prepare a statement
	stmt, err := db_configuration.Prepare("UPDATE configuration SET id=?, name=?, description=?, interval=?;")
	if err != nil {
		return err
	}
	
	// Execute the statement
	_, err = stmt.Exec(n.ID, n.Config.Name, n.Config.Description, n.Config.Interval)

	return err
}

func DatabaseCreateControls()(controls, error){
	var err error
	var Io controls
	
	for i := range Io.ADConvert {
		Io.ADConvert[i]	= converter{false, "", 0, 0}
	}
	
	for i := range Io.DAConvert {
		Io.DAConvert[i]	= converter{false, "", 0, 0}
	}
	
	for i := range Io.DigitalOut {
		Io.DigitalOut[i] = o{false, "", 0}
		ModePin(OutPin[i], out)
	}
	
	for i := range Io.DigitalIO {
		Io.DigitalIO[i] = io{false, "", out, 0}
		ModePin(IOPin[i], Io.DigitalIO[i].Mode)	
	}
	
	for i := range Io.ADConvert {
		_, err = db_controls.Exec("INSERT INTO controls (pinCh, type, enable, description, reference) VALUES($1, $2, $3, $4, $5);",
									ADChannel[i], "ADC", Io.ADConvert[i].Enable, Io.ADConvert[i].Description, Io.ADConvert[i].Reference)
						
	}
	
	for i := range Io.DAConvert {
		_, err = db_controls.Exec("INSERT INTO controls (pinCh, type, enable, description, value, reference) VALUES($1, $2, $3, $4, $5, $6);",
									DAChannel[i], "DAC", Io.DAConvert[i].Enable, Io.DAConvert[i].Description, Io.DAConvert[i].Value, Io.DAConvert[i].Reference)
		
	}
	
	for i := range Io.DigitalOut {
		_, err = db_controls.Exec("INSERT INTO controls (pinCh, type, enable, description, value) VALUES($1, $2, $3, $4, $5);",
									OutPin[i], "O", Io.DigitalOut[i].Enable, Io.DigitalOut[i].Description, Io.DigitalOut[i].Value)
		
	}
	
	for i := range Io.DigitalIO {
		_, err = db_controls.Exec("INSERT INTO controls (pinCh, type, enable, description, value, mode) VALUES($1, $2, $3, $4, $5, $6);",
									IOPin[i], "IO", Io.DigitalIO[i].Enable, Io.DigitalIO[i].Description, Io.DigitalIO[i].Value, Io.DigitalIO[i].Mode)
		
	}

	return Io, err
		
}

func DatabaseGetADC()([nbrAdc]converter, error){
	var err error
	var ADConvert [nbrAdc]converter 	
	
	// Prepare a statement for AD converter
	rows, err := db_controls.Query("SELECT enable, description, reference FROM controls WHERE type = 'ADC';")
	if err != nil {
		return ADConvert, err
	}
	
	var isEmpty bool = false
	// Add all measurements to the slide measurements
	for i := 0; rows.Next() && i < nbrAdc; i++ {
		isEmpty = true
		err = rows.Scan(&ADConvert[i].Enable, &ADConvert[i].Description, &ADConvert[i].Reference);
		if err != nil {
			return ADConvert, err
		}
	}
	
	if !isEmpty {
		err = errors.New("Error: no ADC has already been inserted")
		return ADConvert, err
	}
	
	return ADConvert, err
}

func DatabaseGetDAC()([nbrDac]converter, error){
	var err error
	var DAConvert [nbrDac]converter 	
	
	// Prepare a statement for AD converter
	rows, err := db_controls.Query("SELECT enable, description, value, reference FROM controls WHERE type = 'DAC';")
	if err != nil {
		return DAConvert, err
	}

	var isEmpty bool = false
	// Add all measurements to the slide measurements
	for i := 0; rows.Next() && i < nbrDac; i++ {
		isEmpty = true
		err = rows.Scan(&DAConvert[i].Enable, &DAConvert[i].Description, &DAConvert[i].Value, &DAConvert[i].Reference);
		if err != nil{
			return DAConvert, err
		}
	}
	
	if !isEmpty {
		err = errors.New("Error: no DAC has already been inserted")
		return DAConvert, err
	}
	
	return DAConvert, err
}

func DatabaseGetOutput()([nbrOut]o, error){
	var err error
	var DigitalOut [nbrOut]o    	

	// Prepare a statement for AD converter
	rows, err := db_controls.Query("SELECT enable, description, value FROM controls WHERE type = 'O';")
	if err != nil {
		return DigitalOut, err
	}

	var isEmpty bool = false
	// Add all measurements to the slide measurements
	for i := 0; rows.Next() && i < nbrOut; i++ {
		isEmpty = true
		err = rows.Scan(&DigitalOut[i].Enable, &DigitalOut[i].Description, &DigitalOut[i].Value);
		if err != nil{
			return DigitalOut, err
		}
		ModePin(OutPin[i], out)
	}
	
	if !isEmpty {
		err = errors.New("Error: no Output has already been inserted")
		return DigitalOut, err
	}
	
	return DigitalOut, err
}

func DatabaseGetIO()([nbrIO]io, error){
	var err error
	var DigitalIO [nbrIO]io

	// Prepare a statement for AD converter
	rows, err := db_controls.Query("SELECT enable, description, value, mode FROM controls WHERE type = 'IO';")
	if err != nil {
		return DigitalIO, err
	}

	var isEmpty bool = false
	// Add all measurements to the slide measurements
	for i := 0; rows.Next() && i < nbrIO; i++ {
		isEmpty = true
		err = rows.Scan(&DigitalIO[i].Enable, &DigitalIO[i].Description, &DigitalIO[i].Value, &DigitalIO[i].Mode);
		if err != nil{
			return DigitalIO, err
		}
		ModePin(IOPin[i], DigitalIO[i].Mode)	
	}
	
	if !isEmpty {
		err = errors.New("Error: no IO has already been inserted")
		return DigitalIO, err
	}
	
	return DigitalIO, err
}

func DatabaseGetControls()(controls, error){
	var err error
	var Io controls
	
	// Open a connection with the database
	db_controls, err = sql.Open("sqlite3", controlsDB_filePath)
	if err != nil {
		return Io, err
	}
	// close the database
	defer db_controls.Close()
	
	{
		buffer, err := DatabaseGetADC()
		if err != nil {
			Io, err := DatabaseCreateControls()
			return Io, err
		}
		copy(Io.ADConvert[:], buffer[:])
	}
		
	{
		buffer, err := DatabaseGetDAC()
		if err != nil {
			Io, err := DatabaseCreateControls()
			return Io, err
		}
		copy(Io.DAConvert[:], buffer[:])
	}
	
	{
		buffer, err := DatabaseGetOutput()
		if err != nil {
			Io, err := DatabaseCreateControls()
			return Io, err
		}
		copy(Io.DigitalOut[:], buffer[:])
	}
	
	{
		buffer, err := DatabaseGetIO()
		if err != nil {
			Io, err := DatabaseCreateControls()
			return Io, err
		}
		copy(Io.DigitalIO[:], buffer[:])
	}
	
	return Io, err
	
}

func DatabaseUpdateADC(ADConvert [nbrAdc]converter)(error){
	var err error
	
	// Prepare a statement
	stmt, err := db_controls.Prepare("UPDATE controls SET enable=?, description=?, reference=? WHERE type = 'ADC' AND pinCh = ?;")
	if err != nil {
		return err
	}

	for i := range ADConvert {
		// Execute the statement
		_, err = stmt.Exec(ADConvert[i].Enable, ADConvert[i].Description, ADConvert[i].Reference, ADChannel[i])
	}
		
	return err
}

func DatabaseUpdateDAC(DAConvert [nbrDac]converter)(error){
	var err error
	
	// Prepare a statement
	stmt, err := db_controls.Prepare("UPDATE controls SET enable=?, description=?, value=?, reference=? WHERE type = 'DAC' AND pinCh = ?;")
	if err != nil {
		return err
	}
	
	for i := range DAConvert {
		// Execute the statement
		_, err = stmt.Exec(DAConvert[i].Enable, DAConvert[i].Description, DAConvert[i].Value, DAConvert[i].Reference, DAChannel[i])
	}
		
	return err
}

func DatabaseUpdateOutput(DigitalOut [nbrOut]o)(error){
	var err error
	
	// Prepare a statement
	stmt, err := db_controls.Prepare("UPDATE controls SET enable=?, description=?, value=? WHERE type = 'O' AND pinCh = ?;")
	if err != nil {
		return err
	}
	
	for i := range DigitalOut {
		// Execute the statement
		_, err = stmt.Exec(DigitalOut[i].Enable, DigitalOut[i].Description, DigitalOut[i].Value, OutPin[i])
	}
		
	return err
}

func DatabaseUpdateIO(DigitalIO [nbrIO]io)(error){
	var err error
	
	// Prepare a statement
	stmt, err := db_controls.Prepare("UPDATE controls SET enable=?, description=?, value=?, mode=? WHERE type = 'IO' AND pinCh = ?;")
	if err != nil {
		return err
	}
	
	for i := range DigitalIO {
		// Execute the statement
		_, err = stmt.Exec(DigitalIO[i].Enable, DigitalIO[i].Description, DigitalIO[i].Value, DigitalIO[i].Mode, IOPin[i])
	}
		
	return err
}

func DatabaseUpdateControls(Io controls)(error){
	var err error
	
	// Open a connection with the database
	db_controls, err = sql.Open("sqlite3", controlsDB_filePath)
	if err != nil {
		return err
	}
	// close the database
	defer db_controls.Close()

	err = DatabaseUpdateADC(Io.ADConvert)
	if err != nil {
		return err
	}
	
	err = DatabaseUpdateDAC(Io.DAConvert)
	if err != nil {
		return err
	}
	
	err = DatabaseUpdateOutput(Io.DigitalOut)
	if err != nil {
		return err
	}
	
	err = DatabaseUpdateIO(Io.DigitalIO)

	return err
}
