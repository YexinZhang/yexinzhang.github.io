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

[binlog_row_event_max_size]
将二进制日志中存储的行分组为大小不超过此设置值的事件
如果事件无法拆分，则可以超过最大大小。
选项适用于能够进行基于行的复制的服务器。行以块大小存储在二进制日志中，块大小不超过此选项的值（以字节为单位）。该值必须是256的倍数。默认值为8192。

[MIXED格式]
当以MIXED日志记录格式运行时，服务器在以下情况下自动从基于语句的记录切换为基于行的记录：
When a function contains UUID().


[sync_binlog]
如何刷新binlog至磁盘
1: 同步
0|N: 异步

[binlog_row_event_max_size]
将二进制日志中存储的行分组为大小不超过此设置值的事件
如果事件无法拆分，则可以超过最大大小。


[binlog_encryption]
二进制日志加密
加密自后mysqlbinlog不能直接读取，

[binlog_transaction_compression]
二进制日志压缩， 并不是完全压缩， 有些语句不会

当复制拓扑中的源及其副本都启用了二进制日志事务压缩时,
从主机接收到压缩过的日志，并且将其写入relay_log，解压缩日志并且应用他们，再次压缩写入二进制日志，
下游的从机器会接收到压缩过的日志

当复制拓扑中的源启用了二进制日志事务压缩但其副本没有启用时,
副本接收压缩的事务有效负载并将其压缩后写入其中继日志, 它解压缩事务负载以应用事务，然后将未压缩的事务负载写入其自己的二进制日志（如果有的话）。
任何下游副本均接收未压缩的事务有效负载

如果复制拓扑中的源未启用二进制日志事务压缩，但其副本启用了二进制日志事务压缩，
如果副本具有二进制日志，它将在应用事务负载后压缩它们，并将压缩的事务负载写入其二进制日志。
任何下游副本都将接收压缩的事务有效负载。

当MySQL服务器实例没有二进制日志时，如果它是MySQL 8.0.20的发行版，则无论binlog_transaction_compression的值如何，它都可以接收，处理和显示压缩的事务有效负载。
这样的服务器实例接收到的压缩的事务有效负载以其压缩状态写入中继日志，因此它们间接受益于复制拓扑中其他服务器执行的压缩。

[binary_log_transaction_compression_stats]
perforcemence_schema中的一张表， 用来查看压缩

[events_stages_current]
表
显示事务何时处于其事务有效负载的解压缩或压缩阶段，并显示该阶段的进度。

[binlog_transaction_compression_level_zstd]