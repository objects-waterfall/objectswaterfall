package data

const CreateTable = `CREATE TABLE IF NOT EXISTS %s (
id INTEGER PRIMARY KEY,
json_data TEXT
)`

const CreateWorkerSettingsTable = `CREATE TABLE IF NOT EXISTS Worker_Settings (
    id INTEGER PRIMARY KEY,
    Worker_Name TEXT UNIQUE,
    Timer INTEGER,
    Request_Delay INTEGER,
    Random INTEGER CHECK(Random IN (0,1)),
    Writes_Number_To_Send INTEGER,
    Total_To_Send INTEGER,
    Stop_When_Table_Ends INTEGER CHECK(Stop_When_Table_Ends IN (0,1))
)`

const CreateWorkersResultsTable = `CREATE TABLE IF NOT EXISTS Worker_Results (
	id INTEGER PRIMARY KEY,
    Worker_Id INTEGER,
    Start_Time TEXT,
    Stop_Time TEXT,
    Median_Request_Duration REAL,
    Sended INTEGER,
    Success_Requests INTEGER,
    Failed_Requests INTEGER,
    FOREIGN KEY(Worker_Id) REFERENCES Worker_Settings(id)
)`

const InsertData = `INSERT INTO %s (json_data) values (?)`

const GetJson = `SELECT json_data FROM %s LIMIT ?, ?`

const Count = `Select count () from %s`

// const Tables = `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`
const Workers = `SELECT id, Worker_Name FROM Worker_Settings`
const WorkerName = `SELECT Worker_Name FROM Worker_Settings WHERE id = ?`
const WorkerSettingsId = `SELECT id FROM Worker_Settings WHERE Worker_Name = ?`

const Exists = `SELECT name FROM sqlite_master WHERE type='table' AND name=?`

const InsertWorkerSettings = `INSERT INTO Worker_Settings (
		Worker_Name, 
		Timer, 
		Request_Delay, 
		Random, 
		Writes_Number_To_Send, 
		Total_To_Send, 
		Stop_When_Table_Ends
	) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

const GetWorkerSettings = `SELECT Worker_Name,
		Timer, 
		Request_Delay, 
		Random, 
		Writes_Number_To_Send, 
		Total_To_Send, 
		Stop_When_Table_Ends 
	FROM Worker_Settings 
	WHERE Worker_Name = ?`

const InsertWorkerResults = `INSERT INTO Worker_Results (
	Worker_Id, 
	Start_Time, 
	Stop_Time, 
	Median_Request_Duration, 
	Sended, 
	Success_Requests, 
	Failed_Requests
	)
	VALUES (?, ?, ?, ?, ?, ?, ?);`

const GetWorkerResults = `SELECT
    ws.Worker_Name,
    wr.Start_Time,
    wr.Stop_Time,
    wr.Median_Request_Duration,
    wr.Sended,
    wr.Success_Requests,
    wr.Failed_Requests
	FROM Worker_Settings ws
	JOIN Worker_Results wr ON ws.id = wr.Worker_Id
	WHERE ws.Worker_Name = ?`
