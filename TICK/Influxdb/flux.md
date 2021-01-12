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
  bucket: "qa",
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

# List tag values
import "influxdata/influxdb/schema"
schema.tagValues(bucket: "qa", tag: "host")  # 可以查看机器名

# List tag values in a measurement
# 此函数返回最近30天的结果。
import "influxdata/influxdb/schema"
schema.measurementTagValues(
  bucket: "qa",
  tag: "host",
  measurement: "cpu"
)
```