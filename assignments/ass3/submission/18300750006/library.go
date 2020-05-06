package main

import (
	"fmt"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = ""
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func mustExecute(db *sqlx.DB, SQLs []string) {
	for _, s := range SQLs {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}
type info struct {
	bid string `db:"bid"`
	title string `db:"title"`
	author string `db:"author"`
	ISBN string `db:"ISBN"`
	state string `db:"state"`

}

type Erroccupied string
type Erraccount string
type Errbook string
type Errstate string
type Errhistory string
type ErrDeadline string
type ErrExtend string
type ErrOverdue string
type ErrReturn string
type ErrSuspend string
type ErrArrear string

func (e Erroccupied) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: This bookid '%s' is occupied, please change a new one.\n",string(e))
	return ans
}
func (e Erraccount) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: This student '%s' has already had an account.\n",string(e))
	return ans
}
func (e Errbook) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: No book fits your requirement ")
	return ans
}
func (e Errstate) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: This book %s is cannot be borrowed.",string(e))
	return ans
}
func (e Errhistory) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: No student %s.",string(e))
	return ans
}
func (e ErrDeadline) Error() string{
	var ans string
	ans=fmt.Sprintf("The book %s is not borrowed. \n",string(e))
	return ans
}
func (e ErrExtend) Error() string{
	var ans string
	ans=fmt.Sprintf("You cannot extend the deadline of %s \n",string(e))
	return ans
}
func (e ErrOverdue) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: Fail to check the student %s.",string(e))
	return ans
}
func (e ErrReturn) Error() string{
	var ans string
	ans=fmt.Sprintf("The book %s is not borrowed by this student. \n",string(e))
	return ans
}
func (e ErrSuspend) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: This id %s is invalid \n",string(e))
	return ans
}
func (e ErrArrear) Error() string{
	var ans string
	ans=fmt.Sprintf("Error: This id %s is invalid \n",string(e))
	return ans
}


func (lib *Library) ConnectDB(Userstring, Passwordstring string) error {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", Userstring, Passwordstring, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
	return err
}



// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	mustExecute(lib.db, []string{
		//fmt.Sprintf("DROP DATABASE IF EXISTS %s", DBName),
		//fmt.Sprintf("CREATE DATABASE %s", DBName),
		//fmt.Sprintf("USE %s", DBName),
		"DROP TABLE IF EXISTS borrowed " ,
		"DROP TABLE IF EXISTS record " ,
		"DROP TABLE IF EXISTS book " ,
		"DROP TABLE IF EXISTS student " ,


		//the book in library including its current state
		"CREATE TABLE book (bid CHAR(11) NOT NULL , title CHAR(40)," +
			"author CHAR(40) , ISBN CHAR(20) , state CHAR(10) ," +
			"PRIMARY KEY(bid))",

		//student_info
		"CREATE TABLE student (sname CHAR(20), sid CHAR(11) NOT NULL , sex CHAR(10), " +
			"age int, major CHAR(20)," + "st CHAR(20) ," +
			"PRIMARY KEY(sid))",

		// record of lending
		"CREATE TABLE record (sid CHAR(11) NOT NULL , bid CHAR(11) NOT NULL, " +
			"begintime DATE, endtime DATE," +
			"FOREIGN KEY(sid) references student(sid)," +
			"FOREIGN KEY(bid) references book(bid))",

		//books borrowed now
		"CREATE TABLE borrowed (sid CHAR(11) NOT NULL , bid CHAR(11) NOT NULL, " +
			"begintime DATE, endtime DATE," +
			"extendnum int ," +
			"FOREIGN KEY(sid) references student(sid)," +
			"FOREIGN KEY(bid) references book(bid))",
		// include some books already in the library so that we can test more easily
		"INSERT INTO book VALUES('70826','Othello','W. William Shakespeare','978-0-5215-3517-5','Borrowed')",
		"INSERT INTO book VALUES('71364','Pickwick Papers','Charles Dickens','978-0-1958-6308-6','Available')",
		"INSERT INTO book VALUES('71425','Oliver Twist','Charles Dickens','978-1-8532-6012-4','Available')",
		"INSERT INTO book VALUES('60842','Miser','Molière','978-1-1206-4684-2','Borrowed')",
		"INSERT INTO book VALUES('60843','Miser','Molière','978-1-1206-4684-2','Available')",
		"INSERT INTO book VALUES('91067','Father Goriot','Honoré de Balzac','978-1-4264-0594-5','Available')",
		"INSERT INTO book VALUES('20812','Madame Bovary','Gustave Flauberte','978-1-6740-7054-4','Available')",
		"INSERT INTO book VALUES('94681','One Hundred Years of Solitude','Marquez','978-0-1411-8499-9','Borrowed')",
		"INSERT INTO book VALUES('73498','Crime and punishment','Fyodor Dostoyevsky','978-1-8571-5035-3','Borrowed')",
		"INSERT INTO book VALUES('09728','The Hunger Games','Suzanne Collins','978-0-4390-2348-1','Borrowed')",



		"INSERT INTO student VALUES('hanxy','18300750006','Male','20','CS','Able')",
		"INSERT INTO student VALUES('zhangsan','17307130028','Male','21','CS','Able')",
		"INSERT INTO student VALUES('lisa','16301820009','Female','18','EE','Forbidden')",
		"INSERT INTO borrowed VALUES('18300750006','70826','2020-01-29','2020-05-29',3)",
		"INSERT INTO borrowed VALUES('18300750006','60842','2020-02-29','2020-04-29',1)",
		"INSERT INTO borrowed VALUES('17307130028','94681','2019-12-29','2020-03-29',2)",
		"INSERT INTO borrowed VALUES('17307130028','73498','2020-02-21','2020-04-21',1)",
		"INSERT INTO borrowed VALUES('17307130028','20812','2020-01-26','2020-03-26',1)",
		"INSERT INTO record VALUES('18300750006','71364','2019-03-29','2019-05-29')",
		"INSERT INTO record VALUES('18300750006','91067','2019-07-29','2019-08-29')",

	})
	return nil
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, id, ISBN string) error {
	var bookname string
	//input id, reason
	fmt.Printf("A new book '%s' written by '%s' comes and wants to use the id %s \n",title,author,id)
	errTables := lib.db.QueryRow("SELECT title FROM book WHERE bid = ? ", id).Scan(&bookname)
	if errTables == nil  {
		fmt.Printf("This id is occupied, please input a new one\n")
		return Erroccupied(id)
	}else{
		mustExecute(lib.db, []string{
			fmt.Sprintf("INSERT INTO book VALUES('%s','%s','%s','%s','Available')",
				id , title , author , ISBN),
		})
		fmt.Printf("This book has been successfully added to the library and now is available\n")
		return nil
	}

}

func (lib *Library) RemoveBook(id,reason string) error {

	var bookname string
	//input id, reason
	errTables := lib.db.QueryRow("SELECT title FROM book WHERE bid = ? ", id).Scan(&bookname)
	if  errTables != nil  {
		fmt.Printf("There is no book whose id is '%s' \n",id)
		fmt.Printf("Please input right information.\n")
		return errTables
	}else{
		fmt.Printf("The book '%s' whose id is '%s' is %s \n",bookname,id,reason)
		mustExecute(lib.db, []string{
			fmt.Sprintf("DELETE FROM borrowed WHERE bid = %s",id),
			fmt.Sprintf("UPDATE book SET state = 'Removed' WHERE bid = %s",id),
		})


		fmt.Printf("This book has been successfully removed from the library\n")
		return nil
	}

}

func (lib *Library) AddAccount(sname,sid,sex,age,major,password string) error {

	var saccount string
	var accounttmp string
	//input id, reason

	err1 := lib.db.QueryRow("SELECT DISTINCT user FROM mysql.user WHERE user = 17301257636 ").Scan(&accounttmp)
	if err1 == nil {
		mustExecute(lib.db, []string{
			fmt.Sprintf("DROP USER '17301257636'@'localhost';"),// for convenient test
		})
	}
	errTables := lib.db.QueryRow("SELECT DISTINCT user FROM mysql.user WHERE user = ? ", sid).Scan(&saccount)
	if errTables != nil {
		fmt.Printf("The student '%s' wants to create an account.\n",sid)
		mustExecute(lib.db, []string{
			fmt.Sprintf("CREATE user '%s'@'localhost' identified by '%s';",sid,password),
			fmt.Sprintf("GRANT Insert,Select,Update,Delete ON `ass3`.* TO `%s`@`localhost`;",sid),
			fmt.Sprintf("INSERT INTO student VALUES('%s','%s','%s','%s','%s','Able')",sname,sid,sex,age,major),
		})
		fmt.Printf("Successfully create an account\n")
		return nil
	}else{
		fmt.Printf("The student '%s' wants to create an account\n",sid)
		fmt.Printf("Sorry, one student can just have one account\n")
		return Erraccount(sid)
	}


}

func (lib *Library) QueryBook(manner, value string) error {
	var bid string
	var title string
	var author string
	var ISBN string
	var state string


	if manner == "title" {
		rows,errTables := lib.db.Query("SELECT * FROM book WHERE title = ?",value)
		if errTables != nil  {
			fmt.Printf("No book fits the requirement")
			return Errbook(value)
		}
		fmt.Printf("The followings are the books you query:\n")
		for rows.Next(){
			errTables=rows.Scan(&bid,&title,&author,&ISBN,&state)
			fmt.Printf("%s    %s    %s    %s    %s \n",bid,title,author,ISBN,state)
		}
		fmt.Printf("That's all the books that fit the requirement\n")
		return nil
	}else if manner == "author" {
		rows,errTables := lib.db.Query("SELECT * FROM book WHERE author = ?",value)
		if errTables != nil  {
			fmt.Printf("No book fits the requirement")
			return Errbook(value)
		}
		fmt.Printf("The followings are the books you query:\n")
		for rows.Next(){
			errTables=rows.Scan(&bid,&title,&author,&ISBN,&state)
			fmt.Printf("%s    %s    %s    %s    %s \n",bid,title,author,ISBN,state)

		}
		return nil
	}else if manner == "ISBN" {
		rows,errTables := lib.db.Query("SELECT * FROM book WHERE ISBN = ?",value)
		if errTables != nil {
			return Errbook(value)
		}

		fmt.Printf("The followings are the books you query:\n")
		for rows.Next(){
			errTables=rows.Scan(&bid,&title,&author,&ISBN,&state)
			fmt.Printf("%s    %s    %s    %s    %s \n",bid,title,author,ISBN,state)

		}
		return nil
	}else{
		fmt.Print("Please input the right manner: title,ISBN,author ")
		return Errbook(value)
	}
}

func (lib *Library) BorrowBook(user,bookid string) error {

	var bookstate,studentstate string
	//input id, reason
	errTables := lib.db.QueryRow("SELECT state FROM book WHERE bid = ? ", bookid).Scan(&bookstate)
	lib.db.QueryRow("SELECT st FROM student WHERE sid = ? ", user).Scan(&studentstate)
	//fmt.Printf("%s",studentstate)
	if  errTables != nil {
		//fmt.Println(bookid)
		fmt.Printf("Please input right information.\n")
		fmt.Printf("There is no such book \n")
		return errTables
	}else if  bookstate != "Available" {
		fmt.Printf("The book is not avaliable now\n")
		return Errstate(bookid)
	}else if studentstate != "Able"{
		fmt.Printf("The student can't borrow the book\n")
		//fmt.Printf("%s",studentstate)
		return Errstate(bookid)
	}else{
		mustExecute(lib.db, []string{
			fmt.Sprintf("UPDATE book SET state = 'Borrowed' WHERE bid = %s",bookid),
			fmt.Sprintf("INSERT INTO borrowed VALUES(%s,%s,now(),DATE_ADD(NOW(), interval 1 MONTH),0)",user,bookid),
		})
		fmt.Printf("You have successfully borrowed the book\n")
		return nil
	}

}

func (lib *Library) BorrowHistory(user string) error {
	var sid string
	var bid string
	var begintime string
	var endtime string
	var title string
	var author string
	var ISBN string
	rows,errTables := lib.db.Query("SELECT * FROM record WHERE sid = ? ORDER BY begintime desc",user)
	if  errTables != nil {
		fmt.Printf("Error: Wrong information.\n")
		return Errhistory(user)
	}else{
		fmt.Printf("The followings are the record you query:\n")
		for rows.Next(){
			errTables=rows.Scan(&sid,&bid,&begintime,&endtime)
			err := lib.db.QueryRow("SELECT title, author,ISBN FROM book WHERE bid = ?",bid).Scan(&title,&author,&ISBN)
			if err != nil{
				fmt.Printf("Bug!")
			}else{
				fmt.Printf("%s    %s    %s    %s    %s    %s    %s \n",sid,bid,title,author,ISBN,begintime,endtime)
			}
		}
		fmt.Printf("That's all the record\n")
		return nil
	}
}

func (lib *Library) BorrowNow(user string) error {
	var sid string
	var bid string
	var begintime string
	var endtime string
	var title string
	var author string
	var ISBN string
	rows,errTables := lib.db.Query("SELECT sid,bid,begintime,endtime FROM borrowed WHERE sid = ? ORDER BY begintime desc",user)
	if  errTables != nil {
		fmt.Printf("Error: Wrong information.")
		return Errhistory(user)
	}else{
		fmt.Printf("The followings are the books you borrow now:\n")
		for rows.Next(){
			errTables=rows.Scan(&sid,&bid,&begintime,&endtime)
			err := lib.db.QueryRow("SELECT title, author,ISBN FROM book WHERE bid = ?",bid).Scan(&title,&author,&ISBN)
			if(err != nil){
				fmt.Printf("Bug!")
			}else{
				fmt.Printf("%s    %s    %s    %s    %s    %s    %s \n",sid,bid,title,author,ISBN,begintime,endtime)
			}

		}
		fmt.Printf("That's all the record\n")
		return nil
	}
}

func (lib *Library) CheckDeadline(user,bid string) error {
	var endtime string
	var sid string

	errTables:= lib.db.QueryRow("SELECT sid,endtime FROM borrowed WHERE bid = ?",bid).Scan(&sid,&endtime)
	if  errTables != nil || sid != user {
		fmt.Printf("You don't borrow this book\n")
		return ErrDeadline(bid)
	}else{
		fmt.Printf("The deadline of the book %s is %s \n",bid,endtime)
		fmt.Printf("That's the deadline\n")
		return nil
	}
}

func (lib *Library) ExtendDeadline(user,bid string) error {
	var sid string
	var endtime string
	var extendnum int

	errTables:= lib.db.QueryRow("SELECT sid,endtime,extendnum FROM borrowed WHERE bid = ?",bid).Scan(&sid,&endtime,&extendnum)

	if errTables != nil || sid != user {
		fmt.Printf("You haven't borrowed this book\n")
		return ErrDeadline(bid)
	}else if  extendnum >= 3 {
		fmt.Printf("You cannot extend the deadline anymore becaues you have already extended for 3 times\n")
		return ErrExtend(bid)
	}else{
		mustExecute(lib.db, []string{
			fmt.Sprintf("UPDATE borrowed SET endtime = DATE_ADD(endtime, interval 1 MONTH) WHERE bid = %s",bid),
			fmt.Sprintf("UPDATE borrowed SET extendnum = extendnum + 1 WHERE bid = %s",bid),
		})
		fmt.Printf("Successfully extend\n")
		return nil
	}
}

func (lib *Library) CheckOverdue(user string) error {
	var bid, begintime , endtime,saccount string
	num := 0

	errT := lib.db.QueryRow("SELECT DISTINCT sid FROM student WHERE sid = ? ", user).Scan(&saccount)
	if errT != nil {
		fmt.Printf("No such student, please input the right information\n")
		return ErrOverdue(user)
	}



	rows,errTables := lib.db.Query("SELECT bid,begintime,endtime FROM borrowed WHERE sid = ? AND unix_timestamp(endtime) < unix_timestamp(NOW()) ORDER BY begintime desc",user)
	if  errTables != nil {
		fmt.Printf("Fail to check")
		return ErrOverdue(user)
	}else{
		fmt.Printf("The followings are the overdue books of student %s \n",user)
		for rows.Next(){
			err := rows.Scan(&bid,&begintime,&endtime)
			if err != nil {
				fmt.Printf("Bug!")
			}else{
				fmt.Printf("%s    %s    %s \n",bid,begintime,endtime)
				num = num + 1
			}
		}
		fmt.Printf("That's all the record\n")
		if num != 0  {
			fmt.Printf("You have overdue books! Please return them as soon as possible\n")
		}

		return nil
	}
}

func (lib *Library) ReturnBook(user,bookid string) error {
	var begintime,endtime string
	err := lib.db.QueryRow("SELECT begintime,endtime FROM borrowed WHERE bid = ? AND sid = ? ;",bookid,user).Scan(&begintime,&endtime)

	if err != nil {
		fmt.Printf("You haven't borrowed this book, please check the information you input\n")
		return ErrReturn(bookid)
	}else{
		mustExecute(lib.db, []string{
			fmt.Sprintf("UPDATE book SET state = 'Available' WHERE bid = %s ;",bookid),
			fmt.Sprintf("INSERT INTO record VALUES('%s','%s','%s',now());",user,bookid,begintime),
			fmt.Sprintf("DELETE FROM borrowed WHERE bid = '%s' AND sid = '%s' ;",bookid,user),
		})
		fmt.Printf("Successfully return the book\n")
		return nil
	}
}

func (lib *Library) SuspendAccount(user string) error {
	//DATE_ADD(endtime, interval 1 MONTH)
	var bid ,saccount string
	num := 0
	err := lib.db.QueryRow("SELECT DISTINCT user FROM mysql.user WHERE user = ? ", user).Scan(&saccount)
	if err != nil {
		fmt.Printf("Please input the right userid\n")
		return ErrSuspend(user)
	}

	rows,errTables := lib.db.Query("SELECT bid FROM borrowed WHERE sid = ? AND unix_timestamp(endtime) < unix_timestamp(NOW())",user)
	if errTables != nil {
		return errTables
	}
	for rows.Next(){
		err:=rows.Scan(&bid)
		if err != nil {
			return ErrSuspend(user)
		}else{
			num = num + 1
		}
	}

	if num >= 3 {
		mustExecute(lib.db, []string{
			fmt.Sprintf("UPDATE student SET st = 'Forbidden' WHERE sid = %s",user),
		})
		fmt.Printf("We find that he has many books that are not returned \n")
		fmt.Printf("Successfully suspend the student %s's account \n",user)
	}else{
		fmt.Printf("There is no need to suspend this account\n")
	}

	return nil
}

func (lib *Library) FreeAccount(user string) error{
	var bid ,saccount string
	num := 0
	err := lib.db.QueryRow("SELECT DISTINCT sid FROM student WHERE sid = ? ", user).Scan(&saccount)
	if err != nil {
		fmt.Printf("Please input the right userid\n")
		return ErrSuspend(user)
	}

	rows,errTables := lib.db.Query("SELECT bid FROM borrowed WHERE sid = ? AND unix_timestamp(endtime) < unix_timestamp(NOW())",user)
	if errTables != nil {
		return errTables
	}
	for rows.Next(){
		err:=rows.Scan(&bid)
		if err != nil {
			return ErrSuspend(user)
		}else{
			num = num + 1
		}
	}

	if num == 0 {
		mustExecute(lib.db, []string{
			fmt.Sprintf("UPDATE student SET st = 'Able' WHERE sid = %s",user),
		})
		fmt.Printf("Successfully free the student %s's account \n",user)
	}else{
		fmt.Printf("To free the account, the student should return all the overdue books at first.\n")
	}

	return nil
}

func (lib *Library) Getarrears(user string) error{
	var saccount string
	var sub float64
	var money = 0.0
	err := lib.db.QueryRow("SELECT DISTINCT sid FROM student WHERE sid = ? ", user).Scan(&saccount)
	if err != nil {
		fmt.Printf("Please input the right userid\n")
		fmt.Println(saccount)
		return ErrArrear(user)
	}

	rows,errTables := lib.db.Query("SELECT TIMESTAMPDIFF(DAY,endtime,NOW()) as sub FROM borrowed WHERE sid = ? AND unix_timestamp(endtime) < unix_timestamp(NOW())",user)
	if errTables != nil {
		return errTables
	}
	for rows.Next(){
		err:=rows.Scan(&sub)
		if err != nil {
			return ErrArrear(user)
		}else{
			money = money + sub/10
		}
	}

	fmt.Printf("You need to pay %.1f yuan for the overdue books\n",money)

	return nil
}





func main() {
	fmt.Println("Welcome to the Library Management System!")
	fmt.Println("Please input your account, if you don't have one, please input 'create' to create a new account.")
	var action string
	var sname,sid,sex,age,major,password string
	var loginpassword string
	var wanttodo int
	var title,author,ISBN,bid,reason string
	var mannerid int
	var searchval string
	var loginid string

	// to make the test more conveniently, we will create tables at first, but in reality, there is no need to do this.

	lib2 := Library{}
	lib2.ConnectDB(User,Password)
	lib2.CreateTables()

	flag := 0
	lib := Library{}
	fmt.Scanln(&action)
	if action == "create" {
		fmt.Println("Please input your name.")
		fmt.Scanln(&sname)
		fmt.Println("Please input your student id.")
		fmt.Scanln(&loginid)
		fmt.Println("Please input your password.")
		fmt.Scanln(&password)
		fmt.Println("Please input your sex.")
		fmt.Scanln(&sex)
		fmt.Println("Please input your age.")
		fmt.Scanln(&age)
		fmt.Println("Please input your major.")
		fmt.Scanln(&major)
		lib.ConnectDB(User,Password)
		lib.AddAccount(sname,loginid,sex,age,major,password)
		lib.ConnectDB(loginid,password)
		fmt.Println("Successfully connect to the system")
	}else{
		fmt.Println("Please input your password.")
		fmt.Scanln(&loginpassword)
		loginid = action
		lib.ConnectDB(loginid,loginpassword)
		fmt.Println("Successfully connect to the system")
		if action == User{
			flag = 1
		}
	}
	//fmt.Println(flag)
	for{
		fmt.Println("")
		fmt.Println("What do you want to do?")
		if flag == 1{
			fmt.Println("1. Add a book")
			fmt.Println("2. Remove a book")
			fmt.Println("3. Query the book")
			fmt.Println("4. Query borrow history of a student")
			fmt.Println("5. Query the book being borrowed of a student")
			fmt.Println("6. Query the overdue books of a student")
			fmt.Println("7. Get the arrears of a student")
			fmt.Println("8. Suspend an account")
			fmt.Println("9. Free an account")
			fmt.Println("10. List the table")
			fmt.Println("11. logout")
			fmt.Println("")
		}else{
			fmt.Println("1. Borrow a book")
			fmt.Println("2. Extend the deadline of your book")
			fmt.Println("3. Query the book")
			fmt.Println("4. Query borrow history")
			fmt.Println("5. Query the book being borrowed")
			fmt.Println("6. Query the overdue books")
			fmt.Println("7. Get the arrears")
			fmt.Println("8. Return a book")
			fmt.Println("9. Check the deadline")
			fmt.Println("10. logout")
			fmt.Println("")
		}

		if flag == 1{
			fmt.Scanln(&wanttodo)
			switch wanttodo {
				case 1:
					fmt.Println("Please input the book title")
					fmt.Scanln(&title)
					fmt.Println("Please input the author")
					fmt.Scanln(&author)
					fmt.Println("PLease input the ISBN: xxx-x-xxxx-xxxx-x")
					fmt.Scanln(&ISBN)
					fmt.Println("Please input the bookid you want to give")
					fmt.Scanln(&bid)
					lib.AddBook(title,author,bid,ISBN)
					fmt.Println("Done")
				case 2:
					fmt.Println("Please input the bookid of the book you want to remove")
					fmt.Scanln(&bid)
					fmt.Println("Please input the reason.")
					fmt.Scanln(&reason)
					lib.RemoveBook(bid,reason)
					fmt.Println("Done")
				case 3:
					fmt.Println("Please input the manner you want use to search")
					fmt.Println("1. title")
					fmt.Println("2. author")
					fmt.Println("3. ISBN")
					fmt.Scanln(&mannerid)
					fmt.Println("Please input the value of your search")
					fmt.Scanln(&searchval)
					switch mannerid{
						case 1:
							lib.QueryBook("title",searchval)
						case 2:
							lib.QueryBook("author",searchval)
						case 3:
							lib.QueryBook("ISBN",searchval)
					}
					fmt.Println("Done")
					continuec()

				case 4:
					fmt.Println("Please input the student's id you want to query")
					fmt.Scanln(&sid)
					lib.BorrowHistory(sid)
					fmt.Println("Done")
					continuec()

				case 5:
					fmt.Println("Please input the student's id you want to query")
					fmt.Scanln(&sid)
					lib.BorrowNow(sid)
					fmt.Println("Done")
					continuec()

				case 6:
					fmt.Println("Please input the student's id you want to check")
					fmt.Scanln(&sid)
					lib.CheckOverdue(sid)
					fmt.Println("Done")
					continuec()


				case 7:
					fmt.Println("Please input the student's id you want to search")
					fmt.Scanln(&sid)
					lib.Getarrears(sid)
					fmt.Println("Done")
					continuec()

				case 8:
					fmt.Println("Please input the student's id you want to suspend")
					fmt.Scanln(&sid)
					lib.SuspendAccount(sid)
					fmt.Println("Done")
					continuec()

				case 9:
					fmt.Println("Please input the student's id you want to free")
					fmt.Scanln(&sid)
					lib.FreeAccount(sid)
					fmt.Println("Done")
					continuec()

				case 10:
					fmt.Println("Which table do you want to search?")
					fmt.Println("1. student information")
					fmt.Println("2. book information")
					fmt.Println("3. book being borrowed now")
					fmt.Println("4. borrow record")
					fmt.Scanln(&wanttodo)
					for wanttodo != 1 && wanttodo != 2 && wanttodo != 3 && wanttodo != 4{
						fmt.Println("Please input right number")
						fmt.Scanln(&wanttodo)
					}
					switch wanttodo {
						case 1: lib.List("student")
						case 2: lib.List("book")
						case 3: lib.List("borrowed")
						case 4: lib.List("record")
					}
					fmt.Println("Done")
					continuec()


				case 11:
					fmt.Println("Thank for your use!")
					return
			}
		}else{
			fmt.Scanln(&wanttodo)
			switch wanttodo {
				case 1:
					fmt.Println("Please input the bookid you want to borrow")
					fmt.Scanln(&bid)
					//fmt.Println(loginid)
					lib.BorrowBook(loginid,bid)
					fmt.Println("Done")
					continuec()
				case 2:
					fmt.Println("Please input the bookid of the book whose deadline you want to extend")
					fmt.Scanln(&bid)
					lib.ExtendDeadline(loginid,bid)
					fmt.Println("Done")
					continuec()
				case 3:
					fmt.Println("Please input the manner you want use to search")
					fmt.Println("1. title")
					fmt.Println("2. author")
					fmt.Println("3. ISBN")
					fmt.Scanln(&mannerid)
					fmt.Println("Please input the value of your search")
					fmt.Scanln(&searchval)
					switch mannerid{
					case 1:
						lib.QueryBook("title",searchval)
					case 2:
						lib.QueryBook("author",searchval)
					case 3:
						lib.QueryBook("ISBN",searchval)
					}
					fmt.Println("Done")
					continuec()

				case 4:
					lib.BorrowHistory(loginid)
					fmt.Println("Done")
					continuec()

				case 5:

					lib.BorrowNow(loginid)
					fmt.Println("Done")
					continuec()

				case 6:
					fmt.Println(loginid)
					lib.CheckOverdue(loginid)
					fmt.Println("Done")
					continuec()


				case 7:
					lib.Getarrears(loginid)
					fmt.Println("Done")
					continuec()

				case 8:
					fmt.Println("Please input the bookid of the book you want to return")
					fmt.Scanln(&bid)
					lib.ReturnBook(loginid,bid)
					fmt.Println("Done")
					continuec()

				case 9:
					fmt.Println("Please input the bookid of the book you want to check")
					fmt.Scanln(&bid)
					lib.CheckDeadline(loginid,bid)
					fmt.Println("Done")
					continuec()

				case 10:
					fmt.Println("Thank for your use!")
					return

			}
		}


	}

}

func continuec(){
	var continueid string
	fmt.Println("Press 'c' to continue ")
	fmt.Scanln(&continueid)
	for continueid != "c"{
		fmt.Println("Press 'c' to continue ")
		fmt.Scanln(&continueid)
	}
}

// the following function is just for command line, so i don't add any tests.
//Also, they are so easy that there is no need for us to test them.

func (lib *Library) List(val string) error{
	var bid,title,author,ISBN,state string
	var sname,sid,sex,age,major,st string
	var begintime,endtime,extendnum string

	if val == "book"{
		rows,errTables := lib.db.Query("SELECT * FROM book ")
		if errTables != nil{
			return errTables
		}
		for rows.Next(){
			errTables = rows.Scan(&bid,&title,&author,&ISBN,&state)
			fmt.Printf("%s    %s    %s    %s    %s\n",bid,title,author,ISBN,state)
		}
	} else if val == "student"{
		rows,errTables := lib.db.Query("SELECT * FROM student ")
		if errTables != nil{
			return errTables
		}
		for rows.Next(){
			errTables = rows.Scan(&sname,&sid,&sex,&age,&major,&st)
			fmt.Printf("%s    %s    %s    %s    %s    %s\n",sname,sid,sex,age,major,st)
		}
	} else if val == "borrowed"{
		rows,errTables := lib.db.Query("SELECT * FROM borrowed ")
		if errTables != nil{
			return errTables
		}
		for rows.Next(){
			errTables = rows.Scan(&sid,&sid,&begintime,&endtime,&extendnum)
			fmt.Printf("%s    %s    %s    %s    %s\n",sid,bid,begintime,endtime,extendnum)
		}
	} else if val == "record"{
		rows,errTables := lib.db.Query("SELECT * FROM record ")
		if errTables != nil{
			return errTables
		}
		for rows.Next(){
			errTables = rows.Scan(&sid,&sid,&begintime,&endtime)
			fmt.Printf("%s    %s    %s    %s\n",sid,bid,begintime,endtime)
		}
	}

	return nil

}