# Assignment 1, Environment Setup

This is the first assignment of the course, and it will guide you to setup an environment to be able to play with the actual database system.

The primary database we will be using in this class is MySQL. MySQL is an [open source](https://github.com/mysql/mysql-server) relational database management system created by Oracle, and it is among the topest in the database market[^1] and widely adopted by a lot of companies.

The programming language we will be using in the series of assignments is golang, for its simplicity and the easy-to-use concurrency mechanism. Golang, developed by Google, is a statically typed programming langauge (we'll refer to golang as Go in the following content). It is becomming more and more popular in the recent several years, companies like Bilibili and Zhihu and so many more are using Go to build the backend of their applications, so it is very desirable and meaningful to start to learn this language.

The goal of this assignment is to help you install MySQL on your development environment, and teaches you the basic usage of how to interact with MySQL using command line and programs. In this assignment, you are not required to actually write Go, but you should install the compiler and compile a Go program to interact with the database.

## Environment Reqirements

The operating system is required to be Linux, ideally ubuntu, debian or archlinux, but most of the linux distributions should be fine. We do not have any support for other operating systems like Windows or MacOS, so if you are using any of them there are some suggestions for you:

* Install a virtual machine or dual-boot Linux (mostly recommended)
* Buy a VPS with Linux operating system, and use ssh to access it through out all the assignments

Actually, MacOS and Windows should actually work, because the primary softwares we will be using, MySQl and Go, are all cross platform. However, MySQL as a software that runs in server, it is typically run on Linux, therefore to give you the most native experience of using MySQL, you are highly recommended to try to get access to an Linux operating system for the course.

## Install MySQL and Go

If you think that I am going to teach you how to install them step by step, then you are too naive!

Try to use google to search "install A on B", substitute A with "MySQL" or "golang" and B with the name of your operating system, such as ["install MySQL on Archilinux"](https://lmgtfy.com/?q=install+MySQL+on+Archilinux)[^2]. Note that we will be using MySQL version 8.0, so be careful not to install MySQL 5.7.

Although I am not telling you the steps, there are some common pitfalls that you should be aware of:

* Do not leave the password blank when you are installing MySQL, give it a password for the root user, and memorize it. Later you will use this password as the root user to access the database.

Well, if you are just too lazy, or can not find some good tutorial to teach you how to install them, here is the generous gift from your TA:

* 

[^1]: https://db-engines.com/en/ranking
[^2]: If you do not have access to Google, you should work out a way to get access to it, [here](https://www.uedbox.com/post/54776/) is a website that lists some mirror sites of Google that might be used in mainland China. If you failed to access Google after many attempts, then by all means try to use the international version of [bing](https://bing.com).
