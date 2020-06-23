package main

import (
	"fmt"
	"bufio"
	//"io/ioutil"
	//"path/filepath"
	//"reflect"
	"strings"
	//"sync"
	"time"
	"strconv"
	"os"

	_ "github.com/go-sql-driver/mysql" // this is the driver for mysql
	"github.com/jmoiron/sqlx"      // this is the connector, both package are external packages that you need to use `go install` to install before use
)                        //16

type book struct {
	title string 
	author string
	ISBN string
	//num int
}

/*type borrow struct{
	ISBN string
	stu_id string
	ddl string
}

type add struct{
	ISBN string
	ad_id string
}
                          //35
type remove struct{
	ISBN string
	ad_id string
	explanation string
}*/

var db *sqlx.DB              //42
var maxn = 100
var isbn [100]string
var dif_book_con = 0;
var num [100]int

func yuqi(stid string) int{
	rows,err := db.Query(fmt.Sprintf("select ISBN,ddl from BORROW_RECORD where ddl < now() and stid = '%s'",stid))
	/*stmt,_ := db.Prepare("select ISBN,ddl from BORROW_RECORD where ddl < now() and stid = ?")
	defer stmt.Close()

	rows,err := stmt.Query(stid)*/
	if err != nil {
		fmt.Printf("query error : %v\n",err)
		return 3
	}
	defer rows.Close()
	
	var isbn string
	var ddl string
	fl := 0
	for {
		if !rows.Next() {//end of table
			break
		}

		rows.Scan(&isbn,&ddl)
		ddl_date,_ := time.Parse("2006-01-02", ddl)
  		fmt.Printf("Book %s exceeds ddl: %s\n",isbn,ddl_date)
		fl = 1
	}

	return fl
}

func borrowbook(ISBN string,stid string){
	if(yuqi(stid) == 1){
		fmt.Print("Can't borrow, return immediately.\n")
		return	
	}

	i := 0;
	for i = 0; i < dif_book_con; i++{
		if(isbn[i] == ISBN){break}
	}
	if((i < dif_book_con && num[i] == 0) || i == dif_book_con){
		fmt.Print("Can't borrow, no book.\n")
		return
	}else {
		num[i] --;
		fmt.Printf("Book %s borrowed by student %s\n",ISBN,stid);
	}
	
	stmt,_ := db.Prepare("insert into BORROW_RECORD(ISBN,stid,ddl) values(?,?,DATE_ADD(now(), INTERVAL 3 second))")//DATE_ADD(now(), INTERVAL 30 day)")
	defer stmt.Close()
	
	_,err := stmt.Exec(ISBN,stid)
	if err != nil{
		fmt.Printf("insert data error : %v\n", err)
		return
	}
	return
}

func addbook(ISBN string,adid string){
	
	stmt,_ := db.Prepare("insert into ADD_RECORD(ISBN,adid) values(?,?)")
	defer stmt.Close()
	
	_,err := stmt.Exec(ISBN,adid)
	if err != nil{
		fmt.Printf("insert data error : %v\n", err)
		return
	}
	
	i := 0
	for i = 0; i < dif_book_con; i++{
		if(isbn[i] == ISBN){break}
	}
	if(i < dif_book_con){
		num[i]++
	}
	
	if(i == dif_book_con) {
		fmt.Print("add a new book\n")////
		isbn[dif_book_con] = ISBN
		num[dif_book_con] = 1
		dif_book_con++
		var title,author string
		fmt.Print("title :   ")
		fmt.Scan(&title)
		fmt.Print("\nauthor :  ")
		fmt.Scan(&author)
		fmt.Print("\n")

		stmt,_ := db.Prepare("insert into BOOKS(ISBN,title,author) values(?,?,?)")
		defer stmt.Close()
		_,err := stmt.Exec(ISBN,title,author)
		if err != nil{
			fmt.Printf("insert data error : %v\n", err)
			return
		}
		
	}	
	fmt.Printf("add %d new books\n",1)
}

func addsbook(ISBN string,adid string,nums string){
	
	n, _ := strconv.Atoi(nums)
	stmt,_ := db.Prepare("insert into ADD_RECORD(ISBN,adid) values(?,?)")
	defer stmt.Close()
	
	_,err := stmt.Exec(ISBN,adid)
	if err != nil{
		fmt.Printf("insert data error : %v\n", err)
		return
	}
	
	i := 0
	for i = 0; i < dif_book_con; i++{
		if(isbn[i] == ISBN){break}
	}
	if(i < dif_book_con){
		num[i] += n
		fmt.Printf("add %d new books\n",n)
	}
	
	if(i == dif_book_con) {
		fmt.Print("adding a new book kind\n")////
		isbn[dif_book_con] = ISBN
		num[dif_book_con] = n
		fmt.Printf("add %d new books\n",n)
		dif_book_con++

		var title,author string
		fmt.Print("title :   ")
		fmt.Scan(&title)
		fmt.Print("\nauthor :  ")
		fmt.Scan(&author)
		fmt.Print("\n")

		stmt,_ := db.Prepare("insert into BOOKS(ISBN,title,author) values(?,?,?)")
		defer stmt.Close()
		for j := 0; j < n; j++{
			stmt.Exec(ISBN,title,author)
		}
		
	}	
}

func querybook_ISBN(ISBN string) {
	stmt,_ := db.Prepare("SELECT * FROM BOOKS WHERE ISBN = ?")
	defer stmt.Close()

	rows,err := stmt.Query(ISBN)
	defer rows.Close()
	if err != nil {
		fmt.Printf("query error : %v\n",err)
		return
	}//72

	var temp book
	for rows.Next(){
		rows.Scan(&temp.title,&temp.author,&temp.ISBN)
  		fmt.Println("title:",temp.title,"  author:",temp.author,"  ISBN:",temp.ISBN,"\n")
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func showrecord(){
	var adac,adkey string
	fmt.Print("account :   ")
	fmt.Scan(&adac)
	fmt.Print("\nkeywords :  ")
	fmt.Scan(&adkey)
	fmt.Print("\n")
	//fmt.Printf("%s\n------------\n",adkey)
	stmt,_ := db.Prepare("select keyword from administrator where account = ? ")
	rows,err := stmt.Query(adac)
	defer rows.Close()
	if err != nil {
		fmt.Printf("query error : %v\n",err)
		return
	}
	if rows == nil {
		fmt.Print("Empty account\n")
		return
	}else {

		var key string
		for rows.Next(){
			rows.Scan(&key)
		}
		//fmt.Printf("%s\n------------\n",key)
		if key != adkey{
			fmt.Print("Wrong key\n")
			return
		}
	}

	fmt.Printf("-------------Borrow records------------\n")
	rows,err = db.Query("select ISBN,stid from BORROW_RECORD")
	defer rows.Close()
	if err != nil {
		fmt.Printf("query error : %v\n",err)
		return
	}

	var isbn, stid string
	for rows.Next(){
		rows.Scan(&isbn,&stid)
  		fmt.Println("ISBN:",isbn,"  student:",stid,"\n")
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf(err.Error())
	}
	
	fmt.Printf("-------------ADD records------------\n")
	rows,err = db.Query("select ISBN,adid from ADD_RECORD")
	defer rows.Close()
	if err != nil {
		fmt.Printf("query error : %v\n",err)
		return
	}

	var adid string
	for rows.Next(){
		rows.Scan(&isbn,&adid)
  		fmt.Println("ISBN:",isbn,"  administrator:",stid,"\n")
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func main(){
	db,_ = sqlx.Open("mysql","root:080090@tcp(127.0.0.1:3306)/library")
	if(db == nil){
		fmt.Println("open database fail")//102
		return
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	for{
		fmt.Print("$  ")	
		cmdString,err := reader.ReadString('\n')
		if err != nil{
		fmt.Fprintln(os.Stderr,err)
		}


		//execute
		err = runCommand(cmdString)
		if err != nil{
			fmt.Fprintln(os.Stderr,err)
		}//102
	}
	return
}

func runCommand(commandStr string) error{
	commandStr = strings.TrimSuffix(commandStr,"\n")
	arrCommandStr := strings.Fields(commandStr)
	switch arrCommandStr[0]{
	case "query":
		querybook_ISBN(arrCommandStr[1])
		return nil
	case "add":
		addbook(arrCommandStr[1],arrCommandStr[2])
		return nil	
	case "borrow":
		borrowbook(arrCommandStr[1],arrCommandStr[2])
		return nil
	case "adds":
		addsbook(arrCommandStr[1],arrCommandStr[2],arrCommandStr[3])
		return nil
	
	case "show":
		showrecord()
		return nil
	
	}
return nil
}
