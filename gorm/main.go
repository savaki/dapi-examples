package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/jinzhu/gorm"
	"github.com/savaki/dapi"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Record struct {
	gorm.Model
	Name  string
	Email string
}

func main() {
	var (
		s           = session.Must(session.NewSession(aws.NewConfig()))
		api         = rdsdataservice.New(s)
		driver      = dapi.New(api)
		database    = os.Getenv("DATABASE")
		secretARN   = os.Getenv("SECRET_ARN")
		resourceARN = os.Getenv("RESOURCE_ARN")
		dsn         = fmt.Sprintf("secret=%v resource=%v database=%v", secretARN, resourceARN, database)
		dialect     = "mysql"
	)

	sql.Register(dialect, driver)

	db, err := gorm.Open(dialect, dsn)
	check(err)
	defer db.Close()

	check(db.AutoMigrate(Record{}).Error)

	want := &Record{
		Name:  "name",
		Email: "email",
	}
	err = db.Create(want).Error
	check(err)

	var got Record
	err = db.Model(&Record{}).Where("id = ?", want.ID).First(&got).Error
	check(err)

	got.Name += " updated"
	got.Email += " updated"
	err = db.Model(&Record{}).Update(got).Error
	check(err)

	err = db.Model(&Record{}).Where("id = ?", want.ID).First(&got).Error
	check(err)
}
