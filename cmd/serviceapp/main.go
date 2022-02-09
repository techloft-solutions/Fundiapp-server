package main

import (
	"log"
	"os"

	"github.com/andrwkng/hudumaapp/config"
	"github.com/andrwkng/hudumaapp/database/sqlite"
	"github.com/andrwkng/hudumaapp/server"
	"github.com/go-sql-driver/mysql"
)

var Env string
var Database string

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	log.Println("Config:", cfg)

	//port := os.Getenv("PORT")
	//if port == "" {
	//	log.Fatal("$PORT must be set")
	//}

	//app := NewApp()

	//db := sqlite.NewDB("file:test.db?cache=shared&mode=memory")
	//db := sqlite.NewDB("/hudumaapp.db")
	/*if err := db.Open(); err != nil {
		log.Fatal("cannot open db: %w", err)
	}*/
	var db *sqlite.DB
	server := server.New()
	/*
		dbCfg := mysql.Config{
			User:   "xdshcqjkkzdjs55v",
			Passwd: "whsydeehry48wxsz",
			Net:    "tcp",
			Addr:   "dcrhg4kh56j13bnu.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306",
			DBName: "wtej3mys487jlnyv",
			//Params: nil,
		}*/
	dbCfg := mysql.Config{
		User:   cfg.DB.User,
		Net:    "tcp",
		Addr:   cfg.DB.Addr,
		DBName: cfg.DB.Name,
		Passwd: cfg.DB.Pass,
		Params: nil,
	}

	switch cfg.AppEnv {
	case config.ProdEnv:
		db = sqlite.NewDB(dbCfg.FormatDSN())
		port := os.Getenv("PORT")
		server.Addr = ":" + port
	case config.StageEnv:
		db = sqlite.NewDB(dbCfg.FormatDSN())
		//db = sqlite.NewDB("xdshcqjkkzdjs55v:whsydeehry48wxsz@tcp(dcrhg4kh56j13bnu.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/wtej3mys487jlnyv")
		server.Addr = ":8080"
	default:
		db = sqlite.NewDB("root@tcp(127.0.0.1:3306)/hudumaapp")
		server.Addr = ":8080"
	}

	err = db.Open()
	if err != nil {
		log.Fatal(err)
	}

	server.BkSvc = sqlite.NewBookingService(db)
	server.LocSvc = sqlite.NewLocationService(db)
	server.BidSvc = sqlite.NewBidService(db)
	server.CatSvc = sqlite.NewCategoryService(db)
	server.PfoSvc = sqlite.NewPortfolioService(db)
	server.ReqSvc = sqlite.NewRequestService(db)
	server.UsrSvc = sqlite.NewUserService(db)
	server.RevSvc = sqlite.NewReviewService(db)
	server.IndSvc = sqlite.NewIndustryService(db)
	server.SrchSvc = sqlite.NewSearchService(db)
	server.SubSvc = sqlite.NewSubscriptionService(db)

	log.Fatal(server.Start())

	//_, err := sql.Open("sqlite3", "./hudumaapp.db")*/

}
