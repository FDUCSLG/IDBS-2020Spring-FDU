# Assignment 3, Library

In this assignment, you are going to use what you have learned about MySQL and Go to implement a library management system.

As a sophomore of *Fairy Union of Defense and Attack Nebula* University (FUDAN University) majoring in Computer Science, you are assigned to help build the library management system for the University.

The goal of this assignment to let you know how to design an information management system, and you are given a maximum amount of freedom to design and implement the system, as long as you make use of Go.

## Design

As you are familiar with the library system in Fudan, the requirements are similar to the current library system:

* add books to the library by the administrator's account
* remove a book from the library with explanation (e.g. book is lost)
* add student account by the administrator's account so that the student could borrow books from the library
* query books by title, author or ISBN
* borrow a book from the library with a student account
* query the borrow history of a student account
* query the books a student has borrowed and not returned yet
* check the deadline of returning a borrowed book
* extend the deadline of returning a book, at most 3 times (i.e. refuse to extend if the deadline has been extended for 3 times)
* check if a student has any overdue books that needs to be returned
* return a book to the library by a student account (make sure the student has borrowed the book)
* suspend student's account if the student has more than 3 overdue books (not able to borrow new books unless she has returned books so that she has overdue books less or equal to 3)

And you could add more functionalities to the system as you like.

## Requirement

You should implement the system using Go and mysql. You should abstract the functionalities as Go functions, such as

```go
// AddBook adds a book to the library
func AddBook(book_title string, publisher string) error {
  // query mysql and do the work
}
```

You should document your code with comments and also write tests to prove your implementation for the functionalities is correct.

Take a look at [here](https://blog.alexellis.io/golang-writing-unit-tests/), [here](https://gobyexample.com/testing) or [here](https://labix.org/gocheck) for how to write tests for go programs.

The [simple boilerplate](https://github.com/ichn-hu/IDBS-Spring20-Fudan/tree/master/assignments/ass3/boilerplate) provided might be helpful as a starter, but your are not restricted to it.

Interactions to the system is not required (your tests is already a way of interaction), however it would be a plus if you implement a simple command line interface to interact with the system.

## Submission

You should submit a PDF file as the report explaining how you implement your system and how your implementation fulfills the functionalities. Also include the link to a project on GitHub (create a repository under your own account, and put your project there), and instructions on how to get your system running.

Your report can be in either Chinese or English, and there is no difference in terms of marking, as long as you can make yourself clear.

