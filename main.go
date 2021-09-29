package main

import (
	"database/sql"
	"net/http"
	//"log"
	"fmt"
	"text/template"
	_"github.com/go-sql-driver/mysql"
)
func conectBD()(connect *sql.DB){
  Driver:="mysql"
  User:="root"
  Password:=""
  Name:="System"

  connect,err:=sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
  if err!=nil {
	  panic(err.Error())
  }
  return connect
}
var templates = template.Must(template.ParseGlob("templates/*"))
func main(){
	http.HandleFunc("/", Home)
	http.HandleFunc("/create",Create)
	http.HandleFunc("/insert",Insert)

	http.HandleFunc("/delete",Delete)
	http.HandleFunc("/edit",Edit)

	http.HandleFunc("/update",Update)

	fmt.Println("Start Server..")
	http.ListenAndServe(":8080",nil)
}


type Employee struct {
	Id int
	Name string
	Email string
}
func Home(w http.ResponseWriter, r *http.Request ){
	connectDone := conectBD()
	registry, err := connectDone.Query("SELECT * FROM  employees")

	if  err!=nil {
		panic(err.Error())
	}
	employee:=Employee{}
	arrayEmployee:=[]Employee{}

	for registry.Next(){
		var id int
		var name, email string
		err= registry.Scan(&id,&name,&email)
		if  err!=nil {
			panic(err.Error())
		}
		employee.Id=id
		employee.Name=name
		employee.Email=email

		arrayEmployee=append(arrayEmployee,employee)
		
	}
	templates.ExecuteTemplate(w,"home",arrayEmployee)

}
func Edit(w http.ResponseWriter, r *http.Request ){
	idEmployee:= r.URL.Query().Get("id")
	fmt.Println(idEmployee)
	connectDone := conectBD()
	registry, err := connectDone.Query("SELECT * FROM  employees WHERE id=?",idEmployee)
	
	if  err!=nil {
		panic(err.Error())
	}
	employee:=Employee{}
	for registry.Next(){
		var id int
		var name, email string
		err= registry.Scan(&id,&name,&email)
		if  err!=nil {
			panic(err.Error())
		}
		employee.Id=id
		employee.Name=name
		employee.Email=email
	}
	fmt.Println(employee)
	templates.ExecuteTemplate(w,"edit",employee)
	
}
func Create(w http.ResponseWriter, r *http.Request ){
	//fmt.Fprintf(w, "Hi from CRUD 1")
	templates.ExecuteTemplate(w,"create",nil)

} 
func Update(w http.ResponseWriter, r *http.Request){
	if r.Method=="POST" {
		id:=r.FormValue("id")
		name:=r.FormValue("name")
		email:=r.FormValue("email")
		connectDone := conectBD()
		updateRegistry, err := connectDone.Prepare("UPDATE employees SET name=?,email=? WHERE id=?")
		if  err!=nil {
			panic(err.Error())
		}
		updateRegistry.Exec (name,email,id)
		http.Redirect(w,r,"/",301)
	}
}
func Insert(w http.ResponseWriter, r *http.Request){
	if r.Method=="POST" {
		name:=r.FormValue("name")
		email:=r.FormValue("email")
		connectDone := conectBD()
		insertRegistry, err := connectDone.Prepare("INSERT INTO  employees (name,email) VALUES (?,?)")
		if  err!=nil {
			panic(err.Error())
		}
		insertRegistry.Exec (name,email)
		http.Redirect(w,r,"/",301)
	}
}
func Delete(w http.ResponseWriter, r *http.Request ){
	idEmployee:= r.URL.Query().Get("id")
	fmt.Println(idEmployee)
	connectDone := conectBD()
	deleteRegistry,err := connectDone.Prepare("DELETE FROM employees WHERE id=?")
	if  err!=nil { panic(err.Error())}
	deleteRegistry.Exec(idEmployee)

	http.Redirect(w,r,"/",301)
}