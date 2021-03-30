CREATE USER 'sha2user'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'password';

修改全局变量
SET GLOBAL 
SET PERSIST 修改runtime中的变量，并且将变量永久存在mysql-auto.cnf中
SET PERSIST_ONLY  至保存在文件中，不修改runtime中的变量
RESET PERSIST 删除永久的变量
RESET PERSIST IF EXISTS system_var_name;

SET PERSIST max_connections = 1000;
SET @@PERSIST.max_connections = 1000;

SET PERSIST_ONLY back_log = 100;
SET @@PERSIST_ONLY.back_log = 100;

如果会话系统变量受到限制，则变量描述将指示该限制。 示例包括binlog_format和sql_log_bin。 
设置这些变量的会话值会影响当前会话的二进制日志记录，但也可能对服务器复制和备份的完整性产生更广泛的影响。

结构化变量：
    INSTANCE_NAME.key_buffer_size=xxxx
    其实就是如果我们的业务太多的话，我们可以给一个INSTANCE_NAME来区分变量
    比如
    mtl.key_buffer_size=2M
    qa.key_buffer_size=1M

    default.key_buffer_size=100M 等于原有的 key_buffer_size


Many status variables are reset to 0 by the FLUSH STATUS statement.
FLUSH STATUS可以将大多数状态变量置为零

在函 数 创建中 给 出 DETERMINISTIC 关键 字非 常重要 。 如果 一
个例程对 于 相同的 输入 参数总是 产生相同的 结果, 则认 为该函数
为 DETERMINISTIC ,否则为 NOT DETERMINISTIC 。 如果在例程定义中既未给出DETERMINISTIC ,
也未给出NOT DETERMINISTIC , 则默认是 NOT DETERMINISTIC 。 如果要声
明 一 个函 数是 确定性的 , 则必须明 确指 定 DETERMINISTIC 0
如果将 一 个 NON DET E RMIN I STIC 例程声 明为 DETERMINISTIC ,
可能导致意想不到的 结 果 , 因为它会 导致优 化 器选择不正确 的执行
计 划 。 将 DETERMINISTIC 例程声明 为 NON DETERMINISTIC 可
能会导致可用的优化未被使 用, 降低性能 。


mysql> select date_add(curdate(), interval - 7 day) as '七天前的日期';
+--------------------+
| 七天前的日期         |
+--------------------+
| 2021-03-18         |
+--------------------+

select  from files where file_name like './employees%'\G;
SELECT * FROM INNODB_TABLESPACE\G;
可以查看每个表的信息

INNODB_TABLESTATS
提供近似行数和索引大小

select * from PROCESSLIST\G;
服务器上所有的查询

[innodb_buffer_pool_size]:
InnoDB存储引擎可以使用多少内存空间来缓存内存中的数据和索引
可以动态设置
不要设置的高于数据集的值
单一的mysql主机一般的分配:
RAM | 缓冲池范围
4G --> 1-2G pool size
8 --> 4-6
12 --> 6-10
16 --> 10-12
32 --> 24-28
64 --> 45-56
128 --> 108-116
256 ==> 220-245

[innodb_buffer_pool_instances]
可以将InnoDB缓冲池换分为不同的区域，以便在不同线程读取和写入缓存页面是减少争用，提高并发性。
如果缓冲池大小为12G, instances的值为6, 则每篇区域为2G。

[innodb_log_file_size]
redo log空间大小， 用于数据库崩溃时重放已提交的事务
默认48MB， 测更改需要重启服务器才能生效，开始的时候，可以将其值设置为1G或者2G。

[ACID]
Atomicity
Consistency
Isolation
Duration

START TRANSACTION
BEGIN
SAVEPOINT XXX; 创建回滚点
ROLLBACK;
ROLLBACK TO XXX; 返回回滚点
COMMIT;

隔离级别
show variables like 'transaction_isolation'

读取未提交 read uncommitted
当前事务可以读取由另一个未提交的事务写入的数据，这也称为“脏读”

读提交 read committed
当前事务只能读取另一个事务提交的数据， 这也称为不可重复读

可重复读 repeatable read
一个事务通过第一条语句只能看到相同的数据，即使另一个事务已经提交。在同一个事务中，读取通过第一次读取建立快照是一致的。


SELECT @@SESSION.sql_mode;
SELECT @@GLOBAL.sql_mode;

[sql_mode:] 
    默认值：
        STRICT_TRANS_TABLES: 如果给定的值不能插入事务表中，则中断语句，对于非事务表，只有第一行的给定的值不能插入时才会中断
        NO_ENGINE_SUBSTITUTION: 当诸如CREATE TABLE或ALTER TABLE之类的语句指定被禁用或未在其中编译的存储引擎时，控制默认存储引擎的自动替换。
            启用NO_ENGINE_SUBSTITUTION时，如果所需的引擎不可用，则会发生错误，并且不会创建或更改表。
            在禁用NO_ENGINE_SUBSTITUTION的情况下，对于CREATE TABLE，将使用默认引擎，并且如果所需的引擎不可用，则会发生警告。对于ALTER TABLE，会发生警告，并且该表不会更改。
    
thread_cache_size:
    dynamic
    此值 = 8 + (max_connections / 100)
    客户端链接的线程缓存， (线程池) 当有大量的新链接创建时，调大此值会有比较明显的性能提升
    default: 9

thread_stack:
    static
    每个线程的堆栈大小。默认值足够大，可以正常运行。如果线程堆栈大小太小，则会限制服务器可以处理的SQL语句的复杂性，存储过程的递归深度以及其他消耗内存的操作。

max_connections:
    dynamic
    客户端最大连接数

Connection_errors_max_connections:
    由于达到服务器max_connections限制，拒绝的连接数。


system_time_zone                  | CST               |
time_zone                         | SYSTEM   


[资源组]
resource group
需要在systemd的文件中添加如下信息

```ini
[Service]
AmbientCapabilities=CAP_SYS_NICE
```

MySQL支持创建和管理资源组，并允许将服务器中运行的线程分配给特定的组，以便线程根据该组可用的资源执行。 
使用组属性可以控制其资源，以启用或限制组中线程的资源消耗。 DBA可以根据不同的工作负载修改这些属性。

例如，要管理不需要以高优先级执行的批处理作业的执行，DBA可以创建一个批处理资源组，并根据服务器的繁忙程度上下调整其优先级。 
也许分配给该组的批处理作业应在白天以较低的优先级运行，而在夜间以较高的优先级运行。
DBA还可以调整可用于该组的CPU组。 可以启用或禁用组以控制是否可将线程分配给它们。

The INFORMATION_SCHEMA.RESOURCE_GROUPS table exposes information about resource group definitions and the Performance Schema threads table shows the resource group assignment for each thread.

RESOURCE_GROUPS表显示了资源组，threads表提供了资源组对对应的线程id

每个组都有一个名称。 资源组名称是表和列名称之类的标识符，除非包含特殊字符或保留字，否则不需要在SQL语句中用引号引起来。 
组名不区分大小写，最长不超过64个字符。

每个组都有一个类型，可以是SYSTEM或USER。 资源组类型会影响可分配给该组的优先级值的范围，如下所述。 
此属性以及允许的优先级中的差异使系统线程得以识别，从而保护系统线程免受用户线程争用CPU资源。

系统线程和用户线程与“性能架构线程”表中列出的后台线程和前台线程相对应。

亲和性可以是可用CPU的任何非空子集。 如果组没有亲缘关系，则它可以使用所有可用的CPU。

线程优先级是分配给资源组的线程的执行优先级。 
优先级值的范围是-20（最高优先级）到19（最低优先级）。 
对于系统组和用户组，默认优先级均为0。

System groups are permitted a higher priority than user groups, ensuring that user threads never have a higher priority than system threads:

一般 system的优先级要比user高
    system resource group 的优先级应该被定义为-20-0
    user resource group 应该为0-19
可以启用或禁用每个组，从而使管理员可以控制线程分配。 线程只能分配给已启用的组。

CREATE RESOURCE GROUP 创建
ALTER RESOURCE GROUP  管理
DROP RESOURCE GROUP   删除
Those statements require the RESOURCE_GROUP_ADMIN privilege

SET RESOURCE GROUP    分配

```sql
SELECT * FROM INFORMATION_SCHEMA.RESOURCE_GROUPS\G
CREATE RESOURCE GROUP Batch
  TYPE = USER
  VCPU = 2-3            -- assumes a system with at least 4 CPUs
  THREAD_PRIORITY = 10;
SELECT * FROM INFORMATION_SCHEMA.RESOURCE_GROUPS
       WHERE RESOURCE_GROUP_NAME = 'Batch'\G
SET RESOURCE GROUP Batch FOR thread_id; --- (thread的id可以取threads表中查询)

--- 给当前的链接线程加入到制定的资源组
SET RESOURCE GROUP Batch;

--- To execute a single statement using the Batch group, use the RESOURCE_GROUP optimizer hint:
--- 单个查询语句中使用
INSERT /*+ RESOURCE_GROUP(Batch) */ INTO t2 VALUES(2);

--- 对于系统高负载的情况，减少分配给该组的CPU的数量，降低其优先级，或者（如下所示）：
ALTER RESOURCE GROUP Batch
  VCPU = 3
  THREAD_PRIORITY = 19;

--- 在系统负载较轻的情况下，请增加分配给该组的CPU的数量，提高其优先级，或者（如下所示）：
ALTER RESOURCE GROUP Batch
  VCPU = 0-3
  THREAD_PRIORITY = 0;

```

资源组管理在发生它的服务器上是本地的。 
资源组SQL语句和对resource_groups数据字典表的修改不会写入二进制日志，也不会被复制。


[关闭过程]
```
The shutdown process is initiated.
The server creates a shutdown thread if necessary.
    . If shutdown was requested by a client, a shutdown thread is created.
The server stops accepting new connections.
The server terminates current activity.
The server shuts down or closes storage engines.
    At this stage, the server flushes the table cache and closes all open tables.
每个存储引擎都对其管理的表执行必要的任何操作。 InnoDB将其缓冲池刷新到磁盘（除非innodb_fast_shutdown为2），将当前LSN写入表空间，并终止其自己的内部线程。 MyISAM刷新表的所有挂起的索引写入。
The server exits.
```