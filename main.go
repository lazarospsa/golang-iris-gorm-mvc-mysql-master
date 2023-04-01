package main

import (
	"fmt"
	"gomvc/controllers"
	"gomvc/models"
	"gomvc/repos"
	"gomvc/services"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	// "github.com/kataras/iris/mvc"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	//masterpage
	tmpl := iris.HTML("./templates", ".html").Layout("masterpage.html").Reload(true)
	app.RegisterView(tmpl)

	//static files
	// app.StaticWeb("/static", "./static")

	app.HandleDir("/static", iris.Dir("./static"))

	//routes
	app.Get("/", homeHandler)
	dotenv := goDotEnvVariable("MYSQL_ADDON_HOST")
	fmt.Printf(dotenv)

	// **** BOOKS (MySQL)
	dbhost := os.Getenv("MYSQL_ADDON_HOST")
	dbname := os.Getenv("MYSQL_ADDON_DB")
	dbuser := os.Getenv("MYSQL_ADDON_USER")
	dbpassword := os.Getenv("MYSQL_ADDON_PASSWORD")
	dbport := os.Getenv("MYSQL_ADDON_PORT")

	db, err := gorm.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		app.Logger().Fatalf("error while loading the tables: %v", err)
		return
	}
	//for migrate
	db.AutoMigrate(&models.Book{})

	bookRepo := repos.NewBookRepository(db)
	bookService := services.NewBookService(bookRepo)
	books := mvc.New(app.Party("/books"))
	books.Register(bookService)
	books.Handle(new(controllers.BookController))

	//error
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("error.html")
	})

	//start
	app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
