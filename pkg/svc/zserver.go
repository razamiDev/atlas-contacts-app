package svc

import (
	"fmt"
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"

	pb "github.com/infobloxopen/atlas-contacts-app/pkg/pb"
)

// UniqueViolationCode is a postgres database error
const (
	UniqueViolationCode = "23505"
)

// NewBasicServer returns an instance of the default server interface
func NewBasicServer(database *gorm.DB) (pb.ContactsServer, error) {

	// Callback that allows us to catch create database errors before it gets returned to the response
	database.Callback().Create().After("gorm:create").Register("update_created", updateCreatedAfter)

	return &pb.ContactsDefaultServer{DB: database}, nil
}

// updateCreatedAfter is a callback that allows us to catch database errors and return a human readable error.
func updateCreatedAfter(scope *gorm.Scope) {
	if scope.HasError() {
		err := scope.DB().Error
		if err != nil {
			errPg, ok := err.(*pq.Error) // Change the error to a postgres error
			if ok {
				// Ensuring that we handle unique violation errors from the database and adding our custom message
				if errPg.Code == UniqueViolationCode {
					// Parse the postgres error
					columnName := findColumn(errPg.Detail)
					if columnName == "" {
						columnName = "value"
					}
					valueName := findValue(errPg.Detail)
					var errMsg string
					if valueName == "" {
						errMsg = fmt.Sprintf("A %s already exists with that value", columnName)
					} else {
						errMsg = fmt.Sprintf("A %s already exists with this value (%s)", columnName, valueName)
					}
					scope.DB().Error = fmt.Errorf(errMsg)
					err = scope.DB().Error
				}
			}
		}
	}
}

// Parsing db error helper:
var columnFinder = regexp.MustCompile(`Key \((.+)\)=`)
var valueFinder = regexp.MustCompile(`Key \(.+\)=\((.+)\)`)

// findColumn finds the column in the given pq Detail error string. If the
// column does not exist, the empty string is returned.
func findColumn(detail string) string {
	results := columnFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	} else {
		return results[1]
	}
}

// findColumn finds the column in the given pq Detail error string. If the
// column does not exist, the empty string is returned.
func findValue(detail string) string {
	results := valueFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	} else {
		return results[1]
	}
}

type server struct {
	*pb.ContactsDefaultServer
}
