package tests

import (
	"app/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Seeder struct {
	DB *gorm.DB
}

func NewSeeder() *Seeder {
	gormdb, _ := gorm.Open(mysql.Open(model.GetDSN()), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	return &Seeder{
		DB: gormdb,
	}
}

func (s *Seeder) Seed(fixtures ...interface{}) {
	for _, f := range fixtures {
		s.DB.Create(f)
	}
}

func (s *Seeder) UnSeed(tables ...string) {
	for _, t := range tables {
		s.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", t))
	}
}

func (s *Seeder) Close() {
	db, _ := s.DB.DB()
	db.Close()
}

func (s *Seeder) Dump(w io.Writer, tables ...string) {
	for _, t := range tables {
		switch t {
		case "users":
			d := []model.User{}
			s.DB.Find(&d)
			dumpModelAsJson(w, d)

		case "mail_auths":
			d := []model.MailAuth{}
			s.DB.Find(&d)
			dumpModelAsJson(w, d)
		}
	}
}

func (s *Seeder) DumpStdout(tables ...string) error {
	s.Dump(os.Stdout, tables...)
	return nil
}

func (s *Seeder) DumpFile(filepath string, tables ...string) error {
	if f, err := os.Create(filepath); err != nil {
		return err
	} else {
		s.Dump(f, tables...)
	}
	return nil
}

func dumpModelAsJson(w io.Writer, i interface{}) {
	j, _ := json.MarshalIndent(i, "", "\t")
	r := bytes.NewReader(j)
	io.Copy(w, r)
}
