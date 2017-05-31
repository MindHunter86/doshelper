package apicore

import "time"
import "database/sql"

import "doshelpv2/log"
import "doshelpv2/appctx"

import "golang.org/x/net/context"
import _ "github.com/go-sql-driver/mysql"
import mysql "github.com/go-sql-driver/mysql"


type apiDatabase struct {
	*ApiCore

	db *sql.DB
}
func (self *apiDatabase) configure(ctx context.Context) ( *apiDatabase, error ) {
	if self == nil { return nil,err_glob_InvalidSelf }
	if ctx == nil { return nil,err_glob_InvalidContext }

	self.ApiCore = ctx.Value(appctx.CTX_MOD_APICORE).(*ApiCore)

	if e := self.createConnection(); e != nil { return nil,e }

	self.slogger.W( log.LLEV_DBG, "Database submodule has been initialized and configured!" )
	return self,nil
}
func (self *apiDatabase) createConnection() error {
	if self.db != nil { return err_DB_DdAlreadyDefined }

	var e error
	if self.db, e = sql.Open( "mysql", self.configureConnetcion().FormatDSN() ); e != nil { return e }

	return self.db.Ping()
}
func (self *apiDatabase) destroyConnection() error {
	if self.db == nil { return err_DB_InvalidDb }

	self.slogger.W( log.LLEV_DBG, "Database submodule has been successfully destroyed!" )
	return self.db.Close()
}
func (self *apiDatabase) configureConnetcion() *mysql.Config {
	var cnf *mysql.Config = new(mysql.Config)

	// https://github.com/go-sql-driver/mysql - docs
	cnf.User = self.sql_username
	cnf.Passwd = self.sql_password
	cnf.Net = "tcp4"
	cnf.Addr = self.sql_addr
	cnf.DBName = self.sql_database
	cnf.Collation = "utf8_general_ci"
	cnf.MaxAllowedPacket = 0
	cnf.TLSConfig = "false"
	if tloc, e := time.LoadLocation("Europe/Moscow"); e != nil {	// "Europe%2FMoscow"
		self.slogger.W(log.LLEV_DBG, "Time location parsing error! | " + e.Error())
		cnf.Loc = time.UTC
	} else { cnf.Loc = tloc }

	cnf.Timeout = 10 * time.Second
	cnf.ReadTimeout = 5 * time.Second
	cnf.WriteTimeout = 10 * time.Second

	cnf.AllowAllFiles = false
	cnf.AllowCleartextPasswords = false
	cnf.AllowNativePasswords = false
	cnf.AllowOldPasswords = false
	cnf.ClientFoundRows = false
	cnf.ColumnsWithAlias = false
	cnf.InterpolateParams = false
	cnf.MultiStatements = false
	cnf.ParseTime = true
	cnf.Strict = true // XXX: Only for debug

	return cnf
}

func (self *apiDatabase) getRequestByID(id uint64) ( *request, error ) {
	if self == nil { return nil,err_glob_InvalidSelf }

	return nil,nil
//	stmt, e := self.db.Prepare("select")
}
func (self *apiDatabase) getLastRequest() {}
func (self *apiDatabase) getRequestsCount() {}
func (self *apiDatabase) removeRequestByID(id uint64) {}
