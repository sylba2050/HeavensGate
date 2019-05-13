package main

import (
    "./handler/user"
    "./struct/DB"
    "./utils/myMiddleware"

    "net/http"

    "github.com/labstack/echo"
    _ "github.com/labstack/echo/middleware"

    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db := connectDB("main.sqlite3")
    defer db.Close()

    db.AutoMigrate(&DB.Auth{})
    db.AutoMigrate(&DB.AuthCode{})

    isLoginDB := connectDB("login.sqlite3")
    defer isLoginDB.Close()

    isLoginDB.AutoMigrate(&DB.LoginStatus{})

    e := echo.New()

    e.GET("/code/:userid", user.GenerateAuthCode(db))
    e.POST("/login", user.Login(db, isLoginDB))
    e.POST("/create", user.Create(db))

    g := e.Group("/test", myMiddleware.CustomMiddleware(db))
    g.GET("", test(db))

    e.Start(":8080")
}

func test(db *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.NoContent(http.StatusOK)
    }
}

func connectDB(name string) *gorm.DB {
    db, err := gorm.Open("sqlite3", name)
    if err != nil {
        panic("failed to connect database")
    }
    return db
}
