### Query Result

``` go
type QueryTableResult struct {
	io.Closer
	csvReader     *csv.Reader
	tablePosition int
	tableChanged  bool
	table         *query.FluxTableMetadata
	record        *query.FluxRecord // 返回的Record是这个结构体
	err           error
}

type FluxRecord struct {
	table  int
	values map[string]interface{}
}

host:black,
table:0,
_start:2021-01-12 08:11:29.739731365 +0000 UTC,
_stop:2021-01-12 09:11:29.739731365 +0000 UTC,
cpu:cpu-total,
_measurement:cpu,
result:_result,
_time:2021-01-12 08:11:30 +0000 UTC,
_value:96.26108998819942,
_field:usage_idle

```


### Explore data schema with Flux
```bash
# 函数一般不能混用

# List Buckets
buckets() // List buckets in your organization.

# List measurements
import "influxdata/influxdb/schema"
schema.measurements(bucket: "qa")


# List field keys
import "influxdata/influxdb/schema"
schema.fieldKeys(bucket: "example-bucket")

# List fields in a measurement
schema.measurementFieldKeys(
  bucket: "telegraf",
  measurement: "cpu"
)

# List tag keys
import "influxdata/influxdb/schema"
schema.tagKeys(bucket: "example-bucket")

# List tag keys in a measurement
# 列出在measurement中的tag
schema.measurementTagKeys(
  bucket: "example-bucket",
  measurement: "example-measurement"
)

# 可以查看指定范围内的主机名, 可以快速删掉启用的主机
import "influxdata/influxdb/schema"  
schema.tagValues(
  bucket: "qa",
  tag: "host",
  start: -30s,
)

# List tag values in a measurement
# 此函数返回最近30天的结果。
import "influxdata/influxdb/schema"
schema.measurementTagValues(
  bucket: "qa",
  tag: "host",
  measurement: "cpu"
)
```



### example
```cpp
网卡每秒接受发送的字节数
from(bucket: "qa")
|> range(start: -${simple})
|> filter(fn: (r) => r._measurement == "net" and r.host == "${host}")
|> filter(fn: (r) =>
  r._field == "bytes_recv" or r._field == "bytes_sent"
)
|> aggregateWindow(every: 1m, fn: spread)
|> map(fn: (r) => ({
    r with 
    _value: r._value / 60
}))
```


### 循环
```cpp
if <condition> then <action> else <alternative-action>

if r._value > 95.0000001 and r._value <= 100.0 then "critical"
else if r._value > 85.0000001 and r._value <= 95.0 then "warning"
else if r._value > 70.0000001 and r._value <= 85.0 then "high"
else "normal"
```


### Function definition structure 函数定义
```cpp
// Basic function definition structure
functionName = (functionParameters) => functionOperations

// Function definition
square = (n) => n * n

// Function usage
> square(n:3)

// Function definition
multiply = (x, y) => x * y

// Function usage
> multiply(x:2, y:15)
30
```

#### 管道函数
``` cpp
管道函数
In the example below, the tables parameter is assigned to the <- expression, which represents all data piped-forward into the function. tables is then piped-forward into other operations in the function definition.
在下面的示例中，将表参数分配给<-表达式，该表达式表示所有通过管道传递到函数中的数据。 然后将表通过管道传递到函数定义中的其他操作中。

functionName = (tables=<-) => tables |> functionOperations

// define
下面的示例定义了一个multByX函数，该函数将输入表中每行的_value列乘以x参数。 它使用map（）函数修改每个_value。
// Function definition // tables 只是一个参数名, 可以替换成任意的字符串
multByX = (tables=<-, x) => tables |> map(fn: (r) => ({ r with _value: r._value * x}))
// func usage
// Function usage
from(bucket: "example-bucket")
  |> range(start: -1m)
  |> filter(fn: (r) =>
    r._measurement == "mem" and
    r._field == "used_percent"
  )
  |> multByX(x:2.0)
```

####函数的默认值
Use the = assignment operator to assign a default value to function parameters in your function definition:
```cpp
functionName = (param1=defaultValue1, param2=defaultValue2) => functionOperation

```