### Glossary

```
record influxdb里面插入的每一行记录
```


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
import "influxdata/influxdb/schema"
schema.measurementFieldKeys(
  bucket: "telegraf",
  measurement: "cpu"
)

# List tag keys
import "influxdata/influxdb/schema"
schema.tagKeys(bucket: "example-bucket")

# List tag keys in a measurement
# 列出在measurement中的tag
import "influxdata/influxdb/schema"
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
time format 2019-09-01T00:00:00Z
|> range(start: 2019-09-01T00:00:00Z, stop: 2019-09-01T00:10:00Z)
|> range(start: -$DURATION)

网卡每秒接受发送的字节数
from(bucket: "${buckets}")
|> range(start: -${simple})
|> filter(fn: (r) => r._measurement == "net" and r.host == "${machine}")
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
// 管道函数
In the example below, the tables parameter is assigned to the <- expression, which represents all data piped-forward into the function. tables is then piped-forward into other operations in the function definition.
在下面的示例中，将表参数分配给<-表达式，该表达式表示所有通过管道传递到函数中的数据。 然后将表通过管道传递到函数定义中的其他操作中。

functionName = (tables=<-) => tables |> functionOperations

// define
下面的示例定义了一个multByX函数，该函数将输入表中每行的_value列乘以x参数。 它使用map（）函数修改每个_value。
// Function definition // tables 只是一个参数名, 可以替换成任意的字符串, tables在调用时不用显示指出
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

#### reduce函数
```cpp
函数的定义和管道函数基本类似
functionName = (para,tables) => tables |> reduce(
  identity: {KEY: VALUE, KEY: VALUE...}, //这个identity是自定义的参数，可以供accumulator调用
  fn: (tables, accumulator) => ({
    // 操作 
  })
)

// 函数定义的example:
average = (tabels=<-, outputField="average") => 
  tabels 
    |> reduce(
      identity: {
        count: 1.0,
        sum: 0.0,
        avg: 0.0
      },
      fn: (r, accumulator) => ({
        count: accumulator.count + 1.0,  // 统计record的行数, 每输入一个增加1
        // accumulator.count 的初始值为identity中定义的值， 
        // 计算完成之后结果会存入count中，供下一个record计算调用
        sum: r._value + accumulator.sum, // 计算累加， 每次用r._value和sum相加，并且保留sum
        avg: accumulator.sum / accumulator.count // 通过sum和count的值计算_value的平均数
      }) 
    )
    |> drop(columns: ["count","sum"])
    |> set(key: "_field", value: outputField)
    |> rename(columns: {avg: "_value"})

// 调用example:
from(bucket: "telegraf")
  |> range(start: -10m)
  |> filter(fn: (r) => r._measurement == "mem" and r._field == "used_percent" and r.host == "white")
  |> average()
```

#### 函数的默认值
Use the = assignment operator to assign a default value to function parameters in your function definition:
```cpp
functionName = (param1=defaultValue1, param2=defaultValue2) => functionOperation

```

### stateDuration计算某一个值连续保持了多久
```cpp

// 计算_value为high的值持续了多久，按秒计
from(bucket: "telegraf")
  |> range(start: -10m)
  |> filter(fn: (r) => r._measurement == "mem" and r._field == "used_percent" and r.host == "white")
  |> map(fn: (r) => ({
    r with 
    level: 
      if r._value > 17.99 then "high"
      else "low"
  }))
  |> stateDuration(
    fn: (r) =>
    r.level == "high", 
    unit: 1s)
  |> map(fn: (r) => ({
    r with
    _value: r.stateDuration
  }))
```

### stateCount() 计算连续状态的数量
```cpp
from(bucket: "telegraf")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "mem" and r._field == "used_percent" and r.host == "white")
  |> map(fn: (r) => ({
    r with 
    level: 
      if r._value > 17.99 then "high"
      else "low"
  }))
  |> stateCount(
    fn: (r) => r.level == "high",
  )
  |> map(fn: (r) => ({
    r with 
    _value: r.stateCount  
  })
  )
```