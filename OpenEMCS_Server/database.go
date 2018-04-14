// Server database
package main

import(
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var filePath string = "../Database/nodesDB.db"
var db_nodes *sql.DB

/*															*
* Function/Interface: dataBaseInit
* Param:						
* Return:													
* Description: 	This function create a dataBase if any
*				other exist
*															*/
func DatabaseInit()(error){	
	var err error
	
	// Open a connection with the database
	db_nodes, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}
	// close the database
	defer db_nodes.Close()
	
	// Create the database if any database exists
	_, err = db_nodes.Exec(	"CREATE TABLE IF NOT EXISTS `nodes`(" +
									"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
									"`ip` VARCHAR(15) NULL UNIQUE," + 
									"`name`  VARCHAR(20) NULL," +
									"`description` VARCHAR(100) NULL," +
									"`interval` INTEGER NULL);")
	
	// close the database
	db_nodes.Close()
	
	db_nodes, err = sql.Open("sqlite3", filePath)
	if(err != nil){
		log.Println(err)
		return err
	}
	// close the database
	defer db_nodes.Close()
	
	return err
}

/*															*
* Function/Interface: dataBaseGetNodeIp
* Param:						
	id: node's id which user want to communicate
* Return:						
	string: nodes's ip							
* Description: 	This funtion return an ip address of the node
				indicate by the id
*															*/
func DatabaseGetNodeIP(id int)(string, error){
	var ip string
	var err error

	// Open a connection with the database
	db_nodes, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return ip, err
	}
	// close the database
	defer db_nodes.Close()
	
	// Get ip address at the column indicate by id
	row := db_nodes.QueryRow("SELECT ip FROM nodes WHERE id = $1;", id)
	err = row.Scan(&ip)
	
	return ip, err
}

/*															*
* Function/Interface: dataBaseGetNode
* Param:						
	id: node's id 
* Return:						
	node: a struct node
* Description: 	This funtion return a node with the id and
				the configuration.
*															*/
func DatabaseGetNode(id int)(node, error){
	var err error
	var n node
	
	// Open a connection with the database
	db_nodes, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return n, err
	}
	// close the database
	defer db_nodes.Close()
	
	// Get ip address at the column indicate by id
	row := db_nodes.QueryRow("SELECT name,description,interval FROM nodes WHERE id = $1;", id)
	err = row.Scan(&n.Config.Name, &n.Config.Description, &n.Config.Interval)
	if(err != nil){
		return n, err
	}
	
	n.ID = id
	
	return n, err
}

/*															*
* Function/Interface: dataBaseGetAllNodes
* Param:						
* Return:
	node[]: slide of node with id and configuration struct
* Description: 	This function return all the informations about
				each nodes
*															*/
func DatabaseGetAllNodes()([]node, error){
	var err error
	var nodes []node
	
	// Open a connection with the database
	db_nodes, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return nodes, err
	}
	// close the database
	defer db_nodes.Close()
	
	// Get information about all node into the list
	rows, err := db_nodes.Query("SELECT * FROM nodes;")
	if err != nil {
		return nodes, err
	}
	
	// Add all measurements to the slide measurements
	for rows.Next(){
		var ip string
		var n node
		err = rows.Scan(&n.ID, &ip, &n.Config.Name, &n.Config.Description, &n.Config.Interval);
		if err != nil{
			return nodes, err
		}
		nodes = append(nodes, n)
	}
	
	return nodes, err
}

/*															*
* Function/Interface: databaseSaveNode
* Param:						
	ip: ip address of the new node
	configNode: name and description of the node
* Return:				
	int, node's id
* Description: 	This function create a column for the node
*															*/
func DatabaseSaveNode(ip string, n node)(int, error){
	var err error
	var id int
	var compareIP string
	var stmt *sql.Stmt
	
	// Get the ip of the node if one have already been saved
	if(n.ID != 0) {
		compareIP, _ = DatabaseGetNodeIP(n.ID)
	}
	
	// Open a connection with the database
	db_nodes, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return id, err
	}
	// close the database
	defer db_nodes.Close()
		
	// If it's a new node never used before, or a node with the same id and not same ip already saved
	if(n.ID == 0 || (compareIP != "" && ip != compareIP)) {
		// Prepare a statement, insert if the ip address doesn't exist otherwise update information about node
		stmt, err = db_nodes.Prepare("INSERT OR IGNORE INTO nodes(ip, name, description, interval) VALUES(?, ?, ?, ?);")
		if err != nil {
			return id, err
		}
		
		// Execute the statement
		_, err = stmt.Exec(ip, n.Config.Name, n.Config.Description, n.Config.Interval)
		
	} else {
		// Prepare a statement, insert if the ip address doesn't exist otherwise update information about node
		stmt, err = db_nodes.Prepare("INSERT OR IGNORE INTO nodes(id, ip, name, description, interval) VALUES(?, ?, ?, ?, ?);")
		if err != nil {
			return id, err
		}
		
		// Execute the statement
		_, err = stmt.Exec(n.ID, ip, n.Config.Name, n.Config.Description, n.Config.Interval)
	}
		
	// if the node is already saved into the database
	stmt, err = db_nodes.Prepare("UPDATE nodes SET name=?, description=?, interval=? WHERE ip=?;")
	if err != nil {
		return id, err
	}
	
	// Execute the statement
	_, err = stmt.Exec(n.Config.Name, n.Config.Description, n.Config.Interval, ip)

	// Get id of the new node
	row := db_nodes.QueryRow("SELECT id FROM nodes WHERE ip = $1;", ip)
	err = row.Scan(&id)
	if err != nil {
		return id, err
	}
	
	return id, err
}

/*															*
* Function/Interface: databaseDeleteNode
* Param:						
	ip: ip address of the node that must be deleted
* Return:													
* Description: 	This Delete a node from the database
*															*/
func DatabaseDeleteNode(ip string)(error){
	var err error

	// Open a connection with the database
	db_nodes, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}
	// close the database
	defer db_nodes.Close()
	
	// Delete node at the column indicate by id
	_, err = db_nodes.Exec("DELETE FROM nodes WHERE ip = $1;", ip)
	
	return err
}

/*															*
* Function/Interface: databaseSaveNode
* Param:						
* Return:													
* Description: 	This function delete all the node from the 
				database. It's mainly used at te node startting
*															*/
func DatabaseDeleteAllNodes()(error){
	var err error

	// Open a connection with the database
	db_nodes, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}
	// close the database
	defer db_nodes.Close()

	// Delete all value from database
	_, err = db_nodes.Exec("DELETE FROM nodes;")
	
	return err
}
