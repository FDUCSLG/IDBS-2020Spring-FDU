create table if not exists  ADD_RECORD(
	ISBN varchar(32) not null, 
	adid varchar(32) not null,
	primary key (ISBN,adid)
);
create table if not exists BOOKS(
	ISBN varchar(32) not null primary key,
	title varchar(32) ,
	author varchar(32)
);
create table if not exists BORROW_RECORD(
	ISBN varchar(32) not null,
	stid varchar(32) not null,
	ddl datetime,
	primary key (ISBN,stid)
);
create table if not exists administrator(
	account varchar(32) not null,
	keyword varchar(32) not null
);
insert into administrator(account,keyword) values('001','001');
insert into administrator(account,keyword) values('002','12345');
insert into administrator(account,keyword) values('003','11');
insert into administrator(account,keyword) values('004','sds');
insert into administrator(account,keyword) values('005','sdsaa');

