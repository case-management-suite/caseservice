package main

import (
	"github.com/case-management-suite/caseservice"
	"github.com/case-management-suite/common/config"
)

// main start a gRPC server and waits for connection
func main() {
	appConfig := config.AppConfig{
		CasesService: config.CasesServiceConfig{
			Host: "localhost",
			Port: 7777,
		},
		CasesStorage: config.DatabaseConfig{
			CreateSQL:    "create table if not exists cases(case_id INTEGER PRIMARY KEY, status TEXT);",
			Address:      "./test_cases.db",
			DatabaseType: config.Sqlite,
		},
	}
	caseservice.NewCaseServiceGRPCServer(appConfig).Run()
}
