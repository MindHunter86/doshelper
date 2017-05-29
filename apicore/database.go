package apicore

import "database/sql"

import "doshelpv2/log"
import "doshelpv2/appctx"

import "golang.org/x/net/context"
import _ "github.com/go-sql-driver/mysql"


type apiDatabase struct {
	*ApiCore

	db *sql.DB
}
func (self *apiDatabase) configure(ctx context.Context) ( *apiDatabase, error ) {
	if self == nil { return nil,err_glob_InvalidSelf }
	if ctx == nil { return nil,err_glob_InvalidContext }

	self.ApiCore = ctx.Value(appctx.CTX_MOD_APICORE).(*ApiCore)

	if self.db == nil {
		e := self.createConnection(); if e != nil { return nil,e }
	}

	self.slogger.W( log.LLEV_DBG, "Database submodule has been initialized and configured!" )
	return self,nil
}
func (self *apiDatabase) createConnection() error {
	if self.db != nil { return err_DB_DdAlreadyDefined }

	var e error

	self.db, e = sql.Open("mysql", self.sql_username + ":" + self.sql_password + "@tcp(" + self.sql_host + ":" + self.sql_port + ")/" + self.sql_database)
	if e != nil { return e }

	return self.db.Ping()
}
func (self *apiDatabase) destroyConnection() error {
	if self.db == nil { return err_DB_InvalidDb }

	self.slogger.W( log.LLEV_DBG, "Database submodule has been successfully destroyed!" )
	return self.db.Close()
}
