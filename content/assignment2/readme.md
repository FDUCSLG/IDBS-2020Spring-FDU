# Assignment 2, Evaluator

In the first assignment, you have written several SQLs to help you build basic experience playing with MySQL, in this assignment you are going to continue the journey, but this time, you will use a programming language to interact with MySQL, namely, golang.

Golang, developed by Google, is a statically typed programming langauge (we'll refer to golang as Go in the following content). It is becomming more and more popular in the recent several years, companies like bilibili[^1] [^3] and zhihu[^2] and so many more are using Go to build the backend of their applications. And the idea of cloud native and microservice are built almost on top of Go, just to name a few softwares in Go, docker, kubernetes, and TiDB are built using Go. So it is very desirable and meaningful to start to learn this language.

[^1]: https://github.com/bilibili/kratos
[^2]: https://zhuanlan.zhihu.com/p/48039838
[^3]: https://www.bilibili.com/video/av29079011/

Go is also very simple, it is built with simplicity in mind, as long as you know how to write C code, it won't be hard
to learn Go. However, Go does introduce many ideas that perhaps you are not familiar with, most the the ideas are
related to concurrency, such as `channel` and `goroutine`, we will be using these features extensively in the course,
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
* Try to use GoLand for development, it provides the best Go development environment, and it is free for educational
  purpose.

## Setup

To setup the environment for this assignment, you have the following steps to finish.

### Update Your Cloned Git Repository

Your local repository only have your own code, in order to get the content for this assignment, you need to fetch data
from GitHub. Use the following command:

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


## Concurrent Comparision Result Insert



The single thread version provided in the provided code takes almost 1 minutes on your TA's laptop to insert 7200
records in the the database. This is too slow! The QPS (query per second) is about 7200 / 60 = 120.

Make it concurrent and improve the QPS.

## Evaluate the Score

We evaluate the submission using the following rule. For each SQL file you write, there will be a score of 0 or 1.

For table creation, the rule is, if the table can be created and the `insert_data.sql` could successfully insert data,
then the score is 1, otherwise the score will be 0. The evaluation is already done in the provided code.

For the rest 8 queries, we take a voting approach. For each of the query, we compare every submission's query result
against each other, and if the result is the same, the submission will earn 1 point. For example, when we are evaluating
the submission of `18307130154` on query `1`, we take all submission's result on query `1`, and compare the result of the submission of `18307130154` against
them (including the submission `18307130154` itself), and the count of equal result will be the vote of `18307130154`.
And if submission `18307130154` has the highest vote on query `1`, the score will be 1, otherwise the score is 0.

## Get the Score

