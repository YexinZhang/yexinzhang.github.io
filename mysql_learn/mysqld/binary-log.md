[]

The binary log is resilient(弹性的) to unexpected halts. 
Only complete events or transactions are logged or read back.
Passwords in statements written to the binary log are rewritten by the server not to occur literally in plain text
From MySQL 8.0.14, binary log files and relay log files can be encrypted
A binary log file may become larger than max_binlog_size if you are using large transactions because a transaction is written to the file in one piece, never split between files.
You can delete all binary log files with the RESET MASTER statement, or a subset of them with PURGE BINARY LOGS. 
Within an uncommitted transaction, all updates (UPDATE, DELETE, or INSERT) that change transactional tables such as InnoDB tables are cached until a COMMIT statement is received by the server. 
Modifications to nontransactional tables cannot be rolled back.

[binlog_cache_size]
就是一条事务在未提交之前保存在内存中的二进制日志的大小
在事务期间，用于保存更改的二进制日志的内存缓冲区的大小。该值必须是4096的倍数。
If a statement is bigger than this, the thread opens a temporary file to store the transaction. The temporary file is deleted when the thread ends. 
if binary log encryption is active on the server, the temporary file is encrypted.
请注意，由于复制功能的增强，MySQL 8.0中的二进制日志格式与以前的MySQL版本有所不同。


show status like '%binlog_cache%';
+-----------------------+-------+
| Variable_name         | Value |
+-----------------------+-------+
| Binlog_cache_disk_use | 166   |
| Binlog_cache_use      | 177   |
+-----------------------+-------+
2 rows in set (0.01 sec)


[binlog-error-action]
日志无法记录时的操作
IGNORE_ERROR 如果服务器遇到此类错误，它将继续进行中的事务，记录该错误然后停止记录，并继续执行更新。
ABORT_SERVER 关闭日志记录，并且关闭mysql-server


[binlog_encryption]
二进制日志加密
加密自后mysqlbinlog不能直接读取，


[]