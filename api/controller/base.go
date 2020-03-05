package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nitomibimo/CRUD-GO/api/models"
)

//Server - Import data structure
/*type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

//Initialize - Inititialize of env
func (s *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	e := echo.New()
	db, err := gorm.Open("mysql", "root:@(localhost)/crud_api?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&models.User{}, &models.Post{}) //Auto Migrate from model user and post
	e.Logger.Fatal(e.Start(":8000"))
	db.LogMode(true) // debugging query log

	u := models.User{}
	erQ := db.Model(models.User{}).Scan(&u).Error
	fmt.Println("err query", erQ)

	fmt.Println("hasil ", u)

}*/
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize - connect DB
func (s *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		s.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}
	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		s.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	s.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //database migration

	s.Router = mux.NewRouter()

	s.initializeRoutes()
}

// Run - Listen Port Default
func (s *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
