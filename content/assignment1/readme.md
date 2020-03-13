# Assignment 1, Basic SQL

This is the first assignment of the course, it asks you to write SQL to create tables and query from them.

Basically you need to use

* ceate
* select
* where
* join
* group by
* order by
* having

to be able to finish the assginment.

## Working Directory Setup

This is the first step of the assignment, you need to get yourself familiar with `git` and GitHub.

First of all, create a GitHub id if you do not yet have one.

Then [fork this project](https://github.com/ichn-hu/IDBS-Spring20-Fudan) on GitHub, which will create a mirror repository under your own account.

Then `git clone` your forked project in your local environment. You need to [setup your GitHub account on your local machine] for `git` if you have not yet done so.

Once you've cloned the project in your local environment, get to `IDBS-Spring20-Fudan/assignments/submission` and create a directory and name it by your student id, for example `IDBS-Spring20-Fudan/assignments/submission/16307130177`, and this directory will be your working directory. You are ONLY allowed to modify files under this directory, any other modifications outside this directory will make your submission invalid.

## Table Creation

You need to create 3 tables under the database `ass1`, namely `employee`, `book` and `loan`.

Create a file named `create_table.sql` under your working directory, and write the SQL to create the tables asked in that file.

The constraint for each table is:

* `employee` has 4 columns, `id`, `name`, `office` and `age`, where `id` is the primary key and `name` and `office` is ascii string with length smaller than 32, while `age` is number between 0 and 100. Choose appropriate type for each column.
* `book` also has 4 columns, `id`, `name`, `author` and `publisher`, where `id` is the primary key and other columns are all ascii string with length smaller than 32.
* `record` has 3 columns, `book_id`, `employee_id` and `time`, where the two IDs are foreign keys that references the `id` column of `book` and `employee`, and `time` is of type `DATE`.

Once you finished the table creation SQLs, create a database in your MySQL named `ass1`, and then run the following command in your command line with `username` replaced by your username (such as `root`), the command asks the MySQL client to connect to the MySQL server running at `localhost` on behalf of the user specified by the `username` and run the SQLs in the given file `create_table.sql` under database `ass1`.

```
mysql -h localhost -u username ass1 -p < create_table.sql
```

If `mysql` complaints about any error, try to fix it in `create_table.sql`. You could drop and recreate the database `ass1` if you've ruined the database with wrongly created tables.

Once you've created the table, take a look at `assignments/ass1/insert_data.sql`, which inserts sample data into your newly created tables, run them through MySQL client as well:

```
mysql -h localhost -u username ass1 -p < path/to/assignments/ass1/insert_data.sql
```

## Query from the Database

In this part, you are going to write querys to play with the database you just created.

For each part, create a file in your working directory, for example `1.sql`, `2.sql` etc.

1. Query all fields for employees named `Jones`
2. Query the name of employees with ID equals to `1` or `2` (order dose not matter)
3. Query the name of all employees except the one whose ID is `1` (ordered by ID)
4. Query the name of all employees with age between 25 and 30 (inclusively, ordered by ID)
5. Qeury all fields for employees whose name started with `J` (ordered by age)
6. Query the names of all publishers, if one publisher has more than two books in the database, output the name of the publisher only once (ordered by name, ascii order)
7. Query the id of all boos that is borrowed after `2016-10-31`, also the IDs should be distinct (ordered by id)
8. Query for each employee who has borrowed book more than once, output the `id`, `name`, and number of borrow record (name the field `num`), ordered by `num` in descending order. This one is kind of challenging, the TA's solution uses JOIN, GROUP BY, HAVING and ORDER BY, check them out if you don't know what does these key words mean for MySQL.

## Submit Your Solution

Make sure you have files

* `create_table.sql`
* `1.sql`, ..., `8.sql`

in your working directory.

Then use the following command to submit these files (run it at your working directory).

```
git add .
git commit -m "submission of xxx for ass1"
git push
```

And then get to your GitHub page and create a pull request. The time you create the pull request will be considered as the submission time.

**Note**: the submitted files will be evaluated using a automatic script written in golang, the script will later be uploaded, so make sure you follow all these instructions to make the file hiearchy correct, otherwise the script won't work and you will lose the mark.

Should you have any question, try to search it using Google first. If want clarification of the assignment, then please create an issue in the project repository.
