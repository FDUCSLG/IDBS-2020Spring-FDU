# IDBSass3
进入包含create.sql和library.go的文件夹后，先后输入
mysql -h localhost -u root library -p < create.sql
go run library.go
即可开始访问图书馆数据库
包含的命令有：
query ISBN ：找到对应书本的书籍信息
add ISBN，adid：以管理员身份添加书
adds ISBN，adid，number ：批量增加书
borrow ISBN，stid： 借书
show ：获取历史信息
