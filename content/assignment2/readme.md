# Assignment 2, Evaluator

In the first assignment, you have written several SQLs to help you build basic experience playing with MySQL, in this assignment you are going to continue the journey, but this time, you will use a programming language to interact with MySQL, namely, golang.

Golang, developed by Google, is a statically typed programming langauge (we'll refer to golang as Go in the following content). It is becomming more and more popular in the recent several years, companies like bilibili[^1] [^3] and zhihu[^2] and so many more are using Go to build the backend of their applications. And the idea of cloud native and microservice are built almost on top of Go, just to name a few softwares in Go, docker, kubernetes, and TiDB (a distributed database compatible with MySQL, we probabily will also play with it in the future). So it is very desirable and meaningful to start to learn this language.

[^1]: https://github.com/bilibili/kratos
[^2]: https://zhuanlan.zhihu.com/p/48039838
[^3]: https://www.bilibili.com/video/av29079011/

Go is also very simple, it is built with simplicity in mind, as long as you know how to write C code, it won't be hard
to learn Go. However, Go does introduce many ideas that perhaps you are not familiar with, most of the the ideas are
related to concurrency, such as `channel` and `goroutine`, we will be using these features extensively in the assignments,
since database systems are built to allow (or welcome) concurrent accesses, and Go has a good fame of being a
concurrency friendly programing language.

You might be frightened that you need to learn a new programming language, but it is a must for a good programmer to be
fluent in many programming languages, and it won't be too hard given you already know at least one programming
language. To make your learning experience less painful, follow the instructions below to learn it.

* [Install Go compiler and setup the environment](https://golang.org/doc/install), we will be using the latest Go version, 1.14.
* Join [the tour of Go](https://tour.golang.org/), walk through the whole tutorial, and come back often when you have
  problems writing Go in the future. This is an excellent tutorial, and your TA also learned Go from here. It would take at least 2
  days for you to fully understand the tutorial, so start early. Also, try some of the examples on your local machine,
  which will help you better understand the toolchain of Go. Note, the most important thing for you to learn is the
  concurrency part, you need to understand goroutine (kind of like a light-weighted thread) and `sync.Mutex`, since you
  will need it in this assignment. As for the syntax of Go, it is quite similar with C, so do not worry, you could just
  skim this part and come back when you need it.
* Read the provided code and make sure your understand every line of the code. If you don't, check the tutorial again.
  And if you are still confused, do not hesitate to ask questions in the comment bellow.
* Try to use GoLand for development, it provides the best Go development environment in your TA's opinion, and it is free for educational
  purpose. You could connect GoLand with your MySQL server and GoLand will be able to help you write your SQL.

## Setup

To setup the environment for this assignment, you have the following steps to finish.

### Update Your Cloned Git Repository

Your local repository only have your own code, in order to get the content for this assignment, you need to fetch data
from GitHub. Use the following commands:

```bash
$ git remote add upstream xxx # add your TA's repository to the romote
$ git fetch upstream # this command fetches the update in your TA's repository, it would take a while
$ git merge upstream/master # merge the updates into your master, you can also do `git rebase upstream/master` if you used another branch other than master previously for assignment 1
```

### Copy Provided Code to Your Working Directory

The code provided is under `assignments/ass2/evaluator/`. Similarly you have to create your working directory under
`assignments/ass2/submission/`

Suppose your current directory is the root of the repository.

```bash
$ mkdir assignments/ass2/submission/YOURSTUDENTID # put your student id here
$ cp assignments/ass2/evaluator/* assignments/ass2/submission/YOURSTUDENTID/ # copy provided code into your working directory
```

You can only modify the files in your own working directory, any attempts to modify files outside your working directory
will make your submission invalid.

Especially, only `utils.go` in your working directory will be considered in this assginment (other files will be ignored, see
`assignments/ass2/submission/.gitignore`). With this said, you could actually modify other provided files like
`main.go` if it helps you with understanding the code or debugging your own code, but your modification will not be
considered when we evaluate your submission, and you should make sure your code in `utils.go` can work with the original
codebase.

You should generally only put your code between `YOUR CODE BEGIN` and `YOUR CODE END`. The evaluation of assignment 2 will be based on the three functions you filled in in `utils.go`.

### Create a New Branch for Development

It is your TA's fault that he didn't tell you this in assignment 1. It would be better that you don't work on the master
branch, because that will make the git history messy. Create a new branch and develop your code in the new branch, you
could name the branch `ass2`.

```
$ git checkout -b ass2 # create and checkout to a new branch ass2
$ git add . && git commit -m "update from upstream, start ass2" # commit your modification (only `utils.go` will be tracked) and start working
```

### Modify Configuration

Now edit `utils.go`, in the first `YOUR CODE` block, put your student id in `EvaluateID`, `../../../ass1/submission/`
for `SubmissionDir`, and the user name and password to connect your MySQL server.

You could then run

```bash
$ go run main.go utils.go
```

at your working directory, you should see

```text
submission created
```

in the output if the configuration is correct. If you did not see this line, check the error message and search online, and try to fix it.

You could ignore the printed error message above this line, these messages are for you to know if there
are any problems with the submissions of assignment 1. They are only relevant if you found your solution for assignment 1 is not compliant with
the evaluator, then these error messages tell you why.

If you want to argue about your submission for assignment 1, create another pull request with your modified submission for assignment 1 (only modify files in your working directory for assignment 1) and explain why you think you
deserve a higher score in the pull request, mistakes like creating tables using all UPPERCASE names sound like possible to have a second chance.

Each submission created will have an integer array `score`, and `score[0]` is 0 if the submission's `create_table.sql` is wrong, otherwise it will be 1. `score[1]` to `score[8]` stands for the score for query 1 to query 8 in assignment 1, and they are obtained in the following manner.

## Concurrent Comparision Result Insert

Take a look at `compareAndInsert`, which compares each submissions with all the submissions on each query, and insert a record for the comparison in the table `comparison_result`.

The problem of this function is that it is single-threaded. As you might know, access database is IO-bounded, the CPU might wait for a lot of time for the response of the database, so it makes sense to make the insert concurrent so that the CPU could get more busy and making the insert faster.

The single thread version provided in the provided code takes almost 1 minutes on your TA's laptop to insert 7200
records in the the database. This is too slow! The QPS (query per second) is about 7200 / 60 = 120.

With the help of goroutine, it can be easily made concurrent, so please make it concurrent and improve the QPS in `ConcurrentCompareAndInsert`. As a reference, it takes less than 4 seconds with your TA's simple concurrent version, a speed up of 15x.

You could also do batch insert, which could be even faster, but you won't be able to play with concurrency. If you are interested, you could write a `CompareAndBatchedInsert` function in `utils.go` and test how faster it could be, however this part is unrated.

Also as a reference, your kind TA provided you a concurrent `createSubmissions` in `main.go`, take a look if you don't know how to start.

Also as a kind note, `sql.DB` is not thread-safe, you need multiple `sql.DB` if you want concurrent access of the database, and figuring out (by search it online) what is closure, and how to pass variable to a closure might be helpful.

Run

```bash
$ go run main.go utils.go
```

again, and you will see

```text
ConcurrentCompareAndInsert takes  xxxs
```

instead of

```text
ConcurrentCompareAndInsert Not Implemented
```

## Evaluate the Score

Once you get all the comparision result, it is time to get the score for each submission on each query.

We evaluate the submission using the following rule. In `comparision_result`, for each query, each
submission has record of `is_equal` to every submissions, and the sum of `is_equal` is the vote of
the submission for the query, if the submission has the highest vote, i.e. most of submissions agreed with this submission's result on the query, the score of the submission on this query will be 1, otherwise 0.

You need to use **a single query** to insert the score result in the table `score` created in `createScoreTable`, where `submitter` is the ID of the submitter of the submission, item ranges from 1 to 8 and stands for each query, score being 0 or 1 means the score of `submitter` on query `item`, and `vote` is the vote mentioned above for sanity check.

You need to finish `GetScoreSQL` in `utils.go`, which only returns a single string containing the query sent to MySQL that reads from `comparision_result` and inserts into `score`.

This query can be rather complex and challenging, here are some hints that might be helpful

* UNION is a good wait to combine result of multiple queries
* Common Table Expression (CTE) can be helpful if you want to reuse the result of a subquery
* `order by` can be used with multiple column
* You can insert with values replaced by a select query
* Window functions might be helpful, check it out if you are interested, but you can finish this query without using it

You are free to hack around using any features of MySQL you want (version 8.0), there are many possible ways, as long as you can make it with the `GetScoreSQL` function, you will get the score.

Run

```bash
$ go run main.go utils.go
```

again, and you will see

```text
GetScoreSQL inserted xxx records into score
```

instead of

```text
GetScoreSQL Not Implemented
```

## Get the Score

Finally, you need to fill in each submission's `score` (i.e. `Submission.score`) with the data in table `score`. Finish `GetScore` and you will get output like

```text
...
18307130252 1 1 1 1 1 1 1 1 0
18300200015 1 1 1 1 1 1 1 1 0
18307130031 1 1 1 1 1 1 1 0 1
18307130112 1 1 1 1 1 1 0 1 0
...
```

if you run `go run main.go utils.go`.

## Submit your solution

As you have done for assignment 1, commit your change for `utils.go` and create a pull request to submit it.

Also do not look at other's submission before your submission get merged.

