# Assignment 4, Where to Go from Here

The is the last assignment and it can be fairly easy or extremely hard, depending on how much you like databases.

The first 3 assignments are actually not that hard, they are all created for you to learn and try. The first assignment urges you to write SQL and manually interact with MySQL, and the second one teaches you about how to play with database using a programming language, I especially picked Go for this purpose, for its simplicity and growing attractiveness in the industry, the third one, I have to admit that it is kind of boring, however it asks you to think about how to design the database schema and the API for an application. Overall I believe these assignments are helpful and you must have learnt a lot out of it.

I apologize for such a late release of assignment 4, given the final term is coming and you must have had tons of deadlines, this assignment does not have a due date, and all of you will get a full mark automatically on this assignment. What you need to do for the moment is just to read through the instruction.

## Description

This assignment is a final project, and you are free to do whatever cool with database systems depends on how much you like databases. If you feel that you are not interested in database system, you don't have to do anything for this assignment, or if you don't have time for the moment, you can do it in the summer vacation or whenever you have time.

Since you might not know what to do, I will share with you some resources and inspirations that can help you find what you want to do.

Database system is always an important area in computer science, hence there are a bunch of resources targeting on it. Probably you have already found out by yourself, the following courses are very helpful.

* [15-445/645 Intro to Database Systems (Fall 2019)](https://15445.courses.cs.cmu.edu/fall2019/)  
  Offered by CMU, [has video](https://www.youtube.com/watch?v=oeYBdghaIjc&list=PLSE8ODhjZXjbohkNBWQs_otTrBTrjyohi), introductory course, equivalent to our course, but touches more on the implementation side of a relational database.
* [15-721 Advanced Database Systems (Spring 2019)](https://15721.courses.cs.cmu.edu/spring2020/)  
  Also offered by CMU, [also has video](https://www.youtube.com/watch?v=m72mt4VN9ik&list=PLSE8ODhjZXja7K1hjZ01UTVDnGQdx5v5U), more advanced, research oriented, mainly on in memory database system.

All the course material of the above two courses are freely available, including the assignments and projects, therefore you could look into them to find if you are interested in the projects and try to do it on your own, e.g. implement a concurrent B-Tree in C++.

Other than the above two courses, you might find the [talent plan](https://university.pingcap.com/talent-plan/) of PingCAP attractive to you, including

* [Track 1 TinySQL](https://university.pingcap.com/talent-plan/implement-a-mini-distributed-relational-database/) Implement a distributed relational database in Go
* [Track 2 TinyKV](https://github.com/pingcap-incubator/tinykv) Implement a distributed Key-Value database in Go

Or you can directly get into open source database system development, I strongly suggest you to look at [TiDB](https://github.com/pingcap/tidb), an open source distributed relation database system (disclaimer: I will be employed as a full time Researcher & Developer by PingCAP and help on developing TiDB after graduation). You might find the [track 3](https://university.pingcap.com/talent-plan/become-a-tidb-contributor) and the [source code reading guide](https://pingcap.com/blog-cn/#TiDB-%E6%BA%90%E7%A0%81%E9%98%85%E8%AF%BB) helpful. If you want to work on TiDB, please let me know, and I will be happy to help you get on board.

More over, database system is a fast evolving area, if you want to know what people are doing for future databases, you can look at papers published recently at conferences like [VLDB](https://vldb.org/2019/?papers-research) (stands for Very Large Database) or [SIGMOD](https://sigmod2020.org/sigmod_research_list.shtml) (stands for Special Interest Group on Management of Data).

Above all, I suppose these materials won't get you overflowed but give you a road map of what can be done if you want to dive deeper into database. If you want challenge yourself, please do!

## Submission

You don't have to submit anything, but if you really want to share with us what you have done, send me a PR with an link to the project repository on GitHub.

## Final words

It is indeed a pity that we can not meet each other during the hard time, as I will be graduated soon, perhaps we will never meet, but I really want you to learn well and not get into the same pitfalls I have gotten into. Originally I planned to give some lectures on the implementation side of an actual database system (e.g. TiDB, since I am a contributor to it, and know fairly well about the internals of it) in the classroom, but the pandemic makes it not possible, and I became too busy working on my graduate thesis that I don't have time recording lecture videos. I am really sorry for this.

But I still would like to address some of the logistics beneath my effort of being the TA for you in this introductory database course.

First of all, the assignments, as you can read, are all written in English. It is not mandatory for me to do so, or it is even not encouraged by the instructor, however I still tried so, because being able to read and write in English comfortably is a necessary doorway that leads to a successful career in computer science. I would encourage you to read text books in English instead of translated ones, and try to search in English as well, as most of you are just finishing your first year as a computer science student (considering you only studied "liberal arts" courses as a freshman), it is a good start point to make it a habit using English extensively while studying computer science. I want to tell you this: be comfortable with English, use it as a tool to learn computer science and embrace the world outside of Chinese.

Secondly, I used GitHub as the platform of assignment submission. Many of you did not even have a GitHub account when the first assignment released, and I have helped them setup the account and instructed them how to use git and GitHub. The reason why I used GitHub is also simple, you can not become a programmer without being able to use git and GitHub. Git has become the de facto tool for project version control and Github the primary platform where programmers share code and collaborate. I want you to know that the tool and platform exist, and use them in your future coding life. A piece of advice: maintain an active GitHub account, share your project with the world, and collaborate with people in open source projects as early as possible in your programming career.

Last, I asked you to use Go for the 2 assignments, not only because it is a very convenient language suitable for the assignments (easy concurrency and rich libraries), but also because Go is being increasingly popular in the industry, and we even have database system written in Go (such as [TiDB](https://github.com/pingcap/tidb), [cockroachDB](https://www.cockroachlabs.com/) and [DGraph](https://dgraph.io/)). This is an important lecture I would like to let you know: courses in the university can lag behind the industry and academia really far, try whenever possible new things if you have choice.

I wish you appreciate my efforts and remember the lessons I want to teach you. I am happy to become your TA and friends afterwards. I am actually all around even I have graduated (as I'll work for PingCAP at Shanghai), let me know if you have any questions or difficulties, I am always glad to help.

Wish you happy finals!
