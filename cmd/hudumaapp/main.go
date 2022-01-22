package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andrwkng/hudumaapp/database/sqlite"
	"github.com/andrwkng/hudumaapp/server"
	"github.com/go-sql-driver/mysql"
)

func main() {
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
	cfg := mysql.Config{
		User:   "xdshcqjkkzdjs55v",
		Passwd: "whsydeehry48wxsz",
		Net:    "tcp",
		Addr:   "dcrhg4kh56j13bnu.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306",
		DBName: "wtej3mys487jlnyv",
		Params: nil,
	}

	switch os.Getenv("APP_ENV") {
	case "testing":
		fmt.Println(cfg.FormatDSN())
		db = sqlite.NewDB(cfg.FormatDSN())
		db = sqlite.NewDB("xdshcqjkkzdjs55v:whsydeehry48wxsz@tcp(dcrhg4kh56j13bnu.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/wtej3mys487jlnyv")
		port := os.Getenv("PORT")
		server.Addr = ":" + port
	default:
		// db = sqlite.NewDB("root@tcp(127.0.0.1:3306)/hudumaapp")
		db = sqlite.NewDB(cfg.FormatDSN())
		db = sqlite.NewDB("xdshcqjkkzdjs55v:whsydeehry48wxsz@tcp(dcrhg4kh56j13bnu.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/wtej3mys487jlnyv")
		server.Addr = ":8080"
	}

	err := db.Open()
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

	log.Fatal(server.Start())

	//_, err := sql.Open("sqlite3", "./hudumaapp.db")*/

}
