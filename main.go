package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "github.com/go-sql-driver/mysql"
    "github.com/gin-gonic/gin"
)

var db *sql.DB

type Admin struct {
    Id     	 int  `json:"id"`
    User  	 string  `json:"user"`
    Password string  `json:"password"`
}

type Student struct {
    Id     	 int  `json:"id"`
    Name  	 string  `json:"name"`
    Age		 int  `json:"age"`
}

func main() {
    cfg := mysql.NewConfig()
    cfg.User = "root"
    cfg.Passwd = ""
    cfg.Net = "tcp"
    cfg.Addr = "127.0.0.1:3306"
    cfg.DBName = "dbgolang"

    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

    router := gin.Default()

    router.GET("/admins", getAdmins)
    router.GET("/admins/:id", getAdminByID)
    router.POST("/admins", addAdmin)

    router.GET("/students", getStudents)
    router.GET("/students/:id", getStudentByID)
    router.POST("/students", addStudent)

    router.Run("localhost:8080")
}

func getAdmins(c *gin.Context) {
    var admins []Admin

    rows, err := db.Query("SELECT * FROM admin")
    if err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("admin %q: %v", err))
    }
    defer rows.Close()
    for rows.Next() {
        var alb Admin
        if err := rows.Scan(&alb.Id, &alb.User, &alb.Password); err != nil {
        	c.IndentedJSON(http.StatusOK, fmt.Errorf("admin %q: %v", err))
        }
        admins = append(admins, alb)
    }
    if err := rows.Err(); err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("admin %q: %v", err))
    }
    c.IndentedJSON(http.StatusOK, admins)
}


func getStudents(c *gin.Context) {
    var students []Student

    rows, err := db.Query("SELECT * FROM student")
    if err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("student %q: %v", err))
    }
    defer rows.Close()
    for rows.Next() {
        var alb Student
        if err := rows.Scan(&alb.Id, &alb.Name, &alb.Age); err != nil {
        	c.IndentedJSON(http.StatusOK, fmt.Errorf("student %q: %v", err))
        }
        students = append(students, alb)
    }
    if err := rows.Err(); err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("student %q: %v", err))
    }
    c.IndentedJSON(http.StatusOK, students)
}

func getAdminByID(c *gin.Context) {
	id := c.Param("id")
    var alb Admin

    row := db.QueryRow("SELECT * FROM admin WHERE id = ?", id)
    if err := row.Scan(&alb.Id, &alb.User, &alb.Password); err != nil {
        if err == sql.ErrNoRows {
        	c.IndentedJSON(http.StatusOK, fmt.Errorf("adminById %d: no such admin", id))
        }
        c.IndentedJSON(http.StatusOK, fmt.Errorf("adminById %d: %v", id, err))
    }
    c.IndentedJSON(http.StatusOK, alb)
}

func getStudentByID(c *gin.Context) {
	id := c.Param("id")
    var alb Student

    row := db.QueryRow("SELECT * FROM student WHERE id = ?", id)
    if err := row.Scan(&alb.Id, &alb.Name, &alb.Age); err != nil {
        if err == sql.ErrNoRows {
        	c.IndentedJSON(http.StatusOK, fmt.Errorf("studentById %d: no such student", id))
        }
        c.IndentedJSON(http.StatusOK, fmt.Errorf("studentById %d: %v", id, err))
    }
    c.IndentedJSON(http.StatusOK, alb)
}

func addAdmin(c *gin.Context) {
	var admin Admin
    if err := c.ShouldBindJSON(&admin); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := db.Exec("INSERT INTO admin (user, password) VALUES (?, ?)", admin.User, admin.Password)
    if err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("addAdmin: %v", err))
    }
    id, err := result.LastInsertId()
    if err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("addAdmin: %v", err))
    }
    c.IndentedJSON(http.StatusOK, id)
}

func addStudent(c *gin.Context) {
	var student Student
    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := db.Exec("INSERT INTO student (name, age) VALUES (?, ?)", student.Name, student.Age)
    if err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("addStudent: %v", err))
    }
    id, err := result.LastInsertId()
    if err != nil {
    	c.IndentedJSON(http.StatusOK, fmt.Errorf("addStudent: %v", err))
    }
    c.IndentedJSON(http.StatusOK, id)
}