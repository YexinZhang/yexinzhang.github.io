### db statement struct

<measurement>[,tag-key=tag-val...] field-key=field-val[,field-key=field-val] [unix-nano-timestamp]

    索引在tag上面
    measurement（测量） 相当于sql的库名，field和tag相当于列, 但是tag存储的是元数据，比如hostname,ip,region
    filed 则存储的是测量相关的测量值: 
        比如measurement为cpu_load, field一般会为温度，频率，核心数之类的

    lin protocol
    measurement:  required
    tag set: optional 所有的tag逗号分隔，key和value均为string
    field: Required least one field, key require string type and value can be floats|integers|strings|Booleans.

Retention policy management
    CREATE RETENTION POLICY <retention_policy_name> ON <database> DURATION <duration> REPLICATION <n> [SHARD DURATION <duration>] [default]
    
    DURATION:
        决定数据保存多久 可以设置为持续时间或者INF（infinite），最小一个小时，最大INF
    
    REPLICATION:
        how many independent copies of each point are stored in cluster.
        Replication factors do not serve a purpose with single node instances.
    
    SHARD DURATION
    DEFAULT:
        Set the new retention policy as the default retention policy for the database

    Example:
        create retention policy "jc1d" on mydb duration 1d replication 1  DEFAULT

    Modify with ALTER RETENTION POLICY
        ALTER RETENTION POLICY <retention_policy_name> ON <database_name> DURATION <duration> REPLICATION <n> SHARD DURATION <duration> DEFAULT
        修改的时候必须包含DURATION, REPLICATION, SHARD DURATION, or DEFAULT:中的一个
        example:
            ALTER RETENTION POLICY "jc1d" on "mydb" Duration 1d

    delete retention policy
        DROP RETENTION POLICY <retention_policy_name> ON <database_name>


Create Database

    CREATE DATABASE <database_name> [WITH [DURATION <duration>] [REPLICATION <n>] [SHARD DURATION <duration>] [NAME <retention-policy-name>]]
    

--------------------------------
### Select clause

The basic select statement

Syntax
```
SELECT <field-key>,[<field-key>,<tag-key>] FROM <measurement-name>,[<measurement-name>]
```

SELECT clause
```
    - seletc *
    - select <field-key>
    - select <field-key>,<tag-key> 
        // The select clause must specify at least on field when it include a tag
    - select "<field-key>"::field,"<tag-key>"::tag 
        // ::[field|tag] 用于指明此字段的类型，用于tag和field字段名值相同的情况
        // SELECT "level description"::field,"location"::tag,"water_level"::field FROM "h2o_feet"
        Use this syntax to differentiate between field key and tag key that have the same name
```    

FROM clause
```
    - from <measurement-name>
    - from <measurement-name>,<measurement-name>,...
        单纯查询结果没有聚合，分开打印而已
        ```
        select * from h2o_pH,h2o_feet limit 2
        name: h2o_feet
        time                 level description    location     pH water_level
        ----                 -----------------    --------     -- -----------
        2019-08-17T00:00:00Z below 3 feet         santa_monica    2.064
        2019-08-17T00:00:00Z between 6 and 9 feet coyote_creek    8.12

        name: h2o_pH
        time                 level description location     pH water_level
        ----                 ----------------- --------     -- -----------
        2019-08-17T00:00:00Z                   coyote_creek 7
        2019-08-17T00:00:00Z                   santa_monica 6
        ```
    - FROM <database_name>.<retention_policy_name>.<measurement_name>
    - FROM <database_name>..<measurement_name>
        retention_policy --> default
```
Identifiers must be double quoted if the are contain characters other than [A-z,0-9,_], if the begin with digit, or if they are influxSQL keywords.
尽管并非总是必要，但我们建议您双引号标识符。

Example:
```sql
    select *::field from "h2o_feet" // 查询所有的field
    SELECT ("water_level" * 2) + 4 FROM "h2o_feet"
    SELECT * FROM "NOAA_water_database"."autogen"."h2o_feet" // 完整的from写法，要加上retention_policy
                                            |----> retention_policy
    SELECT * FROM "NOAA_water_database".."h2o_feet" // 默认retention_policy
```
---------------
### Where clause
Syntax
```
    SELECT_clause FROM_clause WHERE <conditional_expression> [(AND|OR) <conditional_expression> [...]]
    where support conditional_expression on fields, tags and timestamp.
```
Note:
```
    where 可以在 field tag timestamp 进行操作
    InfluxDB不支持在WHERE子句中使用OR来指定多个时间范围。
        下面这个语句将返回empty
        > SELECT * FROM "absolutismus" WHERE time = '2016-07-31T20:07:00Z' OR time = '2016-07-31T23:07:17Z'
```
#### Fields contain:
    field_key <operator> ['string' | boolean | float | integer]
    Single quote string field values in WHERE clause.
    Queries with unquoted string field values or double quoted string filed values will not return any data, in most cases, will not return an error.
    字段的值应该被单引号包裹，没有或者被双引号包裹将会返回空，一般情况下，也不会返回错误.
    Support operator:
```
=	equal to
<>	not equal to
!=	not equal to
>	greater than
>=	greater than or equal to
<	less than
<=  less than or equal to
```
#### Tags contain:
```
    tag_key <operator> ['tag-value']
    Single quote tag values in the WHERE clause.
    单引号
    support operator
    =, <>, !=
```

#### Timestamp
    For most SELECT statements, the default time range is between 1677-09-21 00:12:43.145224194 and 2262-04-11T23:47:16.854775806Z UTC.
    For SELECT statements with a GROUP BY time() clause, the default time range is between 1677-09-21 00:12:43.145224194 UTC and now().

Example:
```sql
    > SELECT * FROM "h2o_feet" WHERE "water_level" > 8
    > SELECT * FROM "h2o_feet" WHERE "level description" = 'below 3 feet'
    > SELECT * FROM "h2o_feet" WHERE "water_level" + 2 > 11.9
    > SELECT "water_level" FROM "h2o_feet" WHERE "location" = 'santa_monica'
    > SELECT "water_level" FROM "h2o_feet" WHERE "location" <> 'santa_monica' AND (water_level < -0.59 OR water_level > 9.95)
    > SELECT * FROM "h2o_feet" WHERE time > now() - 7d
        查询过去7天内的数据
    查询时field的string值和tag的值需要带上单引号

        > SELECT "level description" FROM "h2o_feet" WHERE "level description" = at or greater than 9 feet
        ERR: error parsing query: found than, expected ; at line 1, char 86
        > SELECT "level description" FROM "h2o_feet" WHERE "level description" = "at or greater than 9 feet"
        > SELECT "level description" FROM "h2o_feet" WHERE "level description" = 'at or greater than 9 feet'
```
### The Group By clause

    one or more specified tags
    specified time interval
Syntax
```sql
    SELECT_clause FROM_clause [WHERE_clause] GROUP BY [* | <tag_key>[,<tag_key]]
    Example:
        SELECT MEAN("water_level") FROM "h2o_feet" GROUP BY "location"
        SELECT MEAN("water_level") FROM "h2o_feet" GROUP BY *
```
#### Group By Time interval
Basic GROUP BY time() syntax
    SELECT <function>(<field_key>) FROM_clause WHERE <time_range> GROUP BY time(<time_interval>),[tag_key] [fill(<fill_option>)]

Example: // 需要WHERE指定时间段
```sql
SELECT COUNT("water_level") FROM "h2o_feet" WHERE "location"='coyote_creek' AND time >= '2015-08-18T00:00:00Z' AND time <= '2015-08-18T00:30:00Z' GROUP BY time(12m)
```
结果包括了2015-08-18T00:00:00 但是没有包括 2015-08-18T00:12:00, 12被记入下一次的运算



#### 预设时间分析
https://docs.influxdata.com/influxdb/v1.8/query_language/explore-data/#unexpected-timestamps-and-values-in-query-results

原始数据
```sql
> SELECT *  FROM "h2o_feet" WHERE "location"='coyote_creek' AND time >= '2019-08-18T00:06:00Z' AND time < '2019-08-18T00:18:00Z'
```
```
    name: h2o_feet
    time                 level description    location     water_level
    ----                 -----------------    --------     -----------
    2019-08-18T00:06:00Z between 6 and 9 feet coyote_creek 8.419
    2019-08-18T00:12:00Z between 6 and 9 feet coyote_creek 8.32
```

聚合之后的数据(12m)
```sql
SELECT COUNT("water_level") FROM "h2o_feet" WHERE "location"='coyote_creek' AND time >= '2015-08-18T00:06:00Z' AND time < '2015-08-18T00:18:00Z' GROUP BY time(12m)
```
result:
```
    name: h2o_feet
    time                 count
    ----                 -----
    2019-08-18T00:00:00Z 1   //--------注意这里的时间，他是小于指定的值2015-08-18T00:06:00Z
    2019-08-18T00:12:00Z 1
```

InfluxDB对GROUP BY $interval使用预设的整数时间边界，该间隔与WHERE子句中的任何时间条件无关
计算结果时，所有返回的数据必须出现在查询的明确时间范围内，但GROUP BY间隔将基于预设的时间范围。
比如说，这次查询的预设时间范围:
```
    1    time >= 2015-08-18T00:00:00Z' AND time < '2015-08-18T00:12:00Z
    2    time >= 2015-08-18T00:12:00Z' AND time < '2015-08-18T00:24:00Z
```
Group by() interval:
```
    1    time >= 2015-08-18T00:06:00Z AND time < 2015-08-18T00:12:00Z
    2    time >= 2015-08-18T00:12:00Z AND time < 2015-08-18T00:18:00Z
```
Returned Timestamp:"
```
    1    2015-08-18T00:00:00Z
    2    2015-08-18T00:12:00Z
```
import:
    请注意，虽然返回的时间戳记发生在查询的时间范围开始之前，但查询结果排除了查询时间范围之前的数据。

#### Advanced GROUP BY time() Syntax

Syntax
```
SELECT <function>(<field_key>) FROM_clause WHERE <time_range> GROUP BY time(<time_interval>,<offset_interval>),[tag_key] [fill(<fill_option>)]
time(time_interval,offset_interval)
    offset_interval是持续时间文字。它向前或向后移动InfluxDB数据库的预设时间范围。 offset_interval可以为正或负
fill(<fill_option>)
```

offset_interval 用于偏移时间预设值

#### Example1 预设时间段偏移
查询数据集,
```sql
SELECT "water_level" FROM "h2o_feet" WHERE "location"='coyote_creek' AND time >= '2019-08-18T00:00:00Z' AND time <= '2019-08-18T00:54:00Z'

name: h2o_feet
time                 water_level
----                 -----------
2019-08-18T00:00:00Z 8.504
2019-08-18T00:06:00Z 8.419
2019-08-18T00:12:00Z 8.32
2019-08-18T00:18:00Z 8.225
2019-08-18T00:24:00Z 8.13
2019-08-18T00:30:00Z 8.012
2019-08-18T00:36:00Z 7.894
2019-08-18T00:42:00Z 7.772
2019-08-18T00:48:00Z 7.638
2019-08-18T00:54:00Z 7.51
```

#### 不添加offset的查询聚合

注意：查询的时间段从06开始
```sql
SELECT mean("water_level") FROM "h2o_feet" WHERE "location"='coyote_creek' AND time >= '2019-08-18T00:06:00Z' AND time <= '2019-08-18T00:54:00Z' GROUP BY time(18m)

name: h2o_feet
time                 mean
----                 ----
2019-08-18T00:00:00Z 8.3695   
    # 这个结果是 (8.419 + 8.32) / 2 --> 他计算的是06开始到18之间的数据的平均值, 虽然预设值时间段包含了8.504，但是会忽略
2019-08-18T00:18:00Z 8.122333333333334
    # 8.225 8.13 8.012
2019-08-18T00:36:00Z 7.768000000000001
2019-08-18T00:54:00Z 7.51
```
#### 添加offset的查询聚合
```sql
SELECT mean("water_level") FROM "h2o_feet" WHERE "location"='coyote_creek' AND time >= '2019-08-18T00:06:00Z' AND time <= '2019-08-18T00:54:00Z' GROUP BY time(18m,6m)

name: h2o_feet
time                 mean
----                 ----
2019-08-18T00:06:00Z 8.321333333333333   
    # 8.419 8.32 8.225 直接从8.32开始，并且每18m聚合一次数据  
    预设时间起始变为了 06
2019-08-18T00:24:00Z 8.012
2019-08-18T00:42:00Z 7.640000000000001
```

#### 添加负的offset的查询聚合
```sql
SELECT mean("water_level") FROM "h2o_feet" WHERE "location"='coyote_creek' AND time >= '2019-08-18T00:06:00Z' AND time <= '2019-08-18T00:54:00Z' GROUP BY time(18m,-12m)
name: h2o_feet
time                 mean
----                 ----
2019-08-18T00:06:00Z 8.321333333333333   
    # 8.419 8.32 8.225 直接从8.32开始，并且每18m聚合一次数据  
    起始预设时间变成了 2019-08-18T23:48:00Z
    起始预设时间加上间隔18， 时间又到了06， 因为搜索的时间是从06开始的，所以06之前的数据不会被包括，
2019-08-18T00:24:00Z 8.012
2019-08-18T00:42:00Z 7.640000000000001
```


#### GROUP BY time intervals and fill()
主要是更改没有数据的时间间隔报告的值。
Syntax
```sql
SELECT <function>(<field_key>) FROM_clause WHERE <time_range> GROUP BY time(time_interval,[<offset_interval])[,tag_key] [fill(<fill_option>)]
```

fill_opt
``` text
numerical   指定数字
linear      报告没有数据的时间间隔的线性插值结果。
none        Reports no timestamp and no value for time intervals with no data.
null        对于没有数据的时间间隔，报告为null，但返回时间戳。这与默认行为相同。
previous    报告没有数据的时间间隔的前一个时间间隔的值。
```


### INTO clause
INTO子句将查询结果写入用户指定的度量。

Syntax
```sql
SELECT_clause INTO <measurement_name> FROM_clause [WHERE_clause] [GROUP_BY_clause]

- INTO <database_name>.<retention_policy_name>.:MEASUREMENT FROM /<regular_expression>/
    将数据写入用户指定的数据库和保留策略中与FROM子句中的正则表达式匹配的所有measurement。 
    ":MEASUREMENT"是对FROM子句中匹配的每个measurement的反向引用。
        比如FROM "NOAA_water_database"."autogen"."h2o_feet"
        那么 :MEASUREMENT 就等于 h2o_feet
```

example:
```sql
SELECT * INTO "copy_NOAA_water_database"."autogen".:MEASUREMENT FROM "NOAA_water_database"."autogen"./.*/ GROUP BY *
```


### Limit and Slimit
Syntax
```
SELECT_clause [INTO_clause] FROM_clause [WHERE_clause] GROUP BY *[,time(<time_interval>)] [ORDER_BY_clause] LIMIT <N1> SLIMIT <N2>
```
查看measurement的series（序列）, 其实就是measurement的tag的排列组合
```sql
show series from MEASUREMENT
Example:
    show series from h2o_quality
    > show series from h2o_quality
    key
    ---
    h2o_quality,location=coyote_creek,randtag=1
    h2o_quality,location=coyote_creek,randtag=2
    h2o_quality,location=coyote_creek,randtag=3
    h2o_quality,location=santa_monica,randtag=1
    h2o_quality,location=santa_monica,randtag=2
    h2o_quality,location=santa_monica,randtag=3    
```
slimit 必须带上group by *字段才能正常使用

```
    limit N --> 返回N个数据点
    slimit N --> 返回N个序列的数据点
    limit X slimit Y --> 返回前Y个序列的X个数据点
```


### Time zone clause
The tz() clause returns the UTC offsets for the specified timezone

Syntax
```
SELECT_clause [INTO_clause] FROM_clause [WHERE_clause] [GROUP_BY_clause] [ORDER_BY_clause] [LIMIT_clause] [OFFSET_clause] [SLIMIT_clause] [SOFFSET_clause] tz('<time_zone>')
```

By default, influxDB returns and stores timestamps in UTC. The tz() clause includes the UTC offset.