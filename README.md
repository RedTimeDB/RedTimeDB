<!--
 * @Author: gitsrc
 * @Date: 2022-04-02 14:21:40
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-18 17:02:54
 * @FilePath: /RedEpochDB/README.md
-->

<p align="center">
<img 
    src="./imgs/logo.png" width="80%" height="80%"
    border="0" alt="RedEpochDB" />
</p>


# RedEpochDB

![build](https://github.com/RedEpochDB/RedEpochDB/actions/workflows/build.yml/badge.svg) 
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FRedEpochDB%2FRedEpochDB.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FRedEpochDB%2FRedEpochDB?ref=badge_shield)

The next-generation time series database, builds ultra-high-speed write and query performance based on the memory model, and has remote persistent write support,Welcome to experience a faster time series database.

1. Ultra-high data write performance.
2. Excellent remote data persistence support.
3. Support cluster horizontal expansion.
4. Support community popular chart system.
5. Support multiple protocols for write access: MQTT, HTTP, RESP.

# Architecture
<p align="center">
<img 
    src="https://github.com/user-attachments/assets/f6c1ab29-5208-47c6-81c7-9d7b3bf77156" 
     alt="Architecture">
</p>

# Model test
## metrics throughput
<p align="center">
<img 
    src="./imgs/test_metrics.jpg" 
    border="0" alt="test_metrics" width="80%" />
</p>

## rows throughput
<p align="center">
<img 
    src="./imgs/test_rows.jpg" 
    border="0" alt="test_rows" width="80%" />
</p>

# Command support
## TS.ADD
### Instruction description
Append a sample to a time series.

### Syntax
```
TS.ADD key timestamp value 
  [RETENTION retentionPeriod] 
  [ENCODING [COMPRESSED|UNCOMPRESSED]] 
  [CHUNK_SIZE size] 
  [ON_DUPLICATE policy] 
  [LABELS {label value}...]
```

### Example of use
```
127.0.0.1:6380> TS.ADD temperature:3:11 1651708850000 27
(integer) 1651708850000
127.0.0.1:6380> TS.ADD temperature:3:11 1651708860000 28
(integer) 1651708860000
127.0.0.1:6380> TS.ADD temperature:3:11 1651708870000 29
(integer) 1651708870000
```

## TS.GET
### Instruction description
Get the last sample.
### Syntax
```
TS.GET key 
  [LATEST]
```
### Example of use
```
127.0.0.1:6380> TS.GET temperature:3:11
1) (integer) 1651708850000
2) "27"
```

## TS.MADD
### Instruction description
Append new samples to one or more time series.

### Syntax
```TS.MADD {key timestamp value}...```

### Example of use
```
127.0.0.1:6380> TS.MADD cpu_usage_user{1687509904} 1651708850000 33 cpu_usage_system{1687509904} 1651708850000 70 cpu_usage_idle{1687509904} 1651708850000 79 cpu_usage_nice{1687509904} 1651708850000 6 cpu_usage_iowait{1687509904} 1651708850000 56 cpu_usage_irq{1687509904} 1651708850000 29 cpu_usage_softirq{1687509904} 1651708850000 65 cpu_usage_steal{1687509904} 1651708850000 63 cpu_usage_guest{1687509904} 1651708850000 63 cpu_usage_guest_nice{1687509904} 1651708850000 83
 1) "1651708850000"
 2) "1651708850000"
 3) "1651708850000"
 4) "1651708850000"
 5) "1651708850000"
 6) "1651708850000"
 7) "1651708850000"
 8) "1651708850000"
 9) "1651708850000"
10) "1651708850000"
```

## TS.RANGE
### Instruction description
Query a range in forward direction.
### Syntax
```
TS.RANGE key fromTimestamp toTimestamp
  [LATEST]
  [FILTER_BY_TS ts...]
  [FILTER_BY_VALUE min max]
  [COUNT count] 
  [[ALIGN align] AGGREGATION aggregator bucketDuration [BUCKETTIMESTAMP bt] [EMPTY]]
```
### Example of use
```
127.0.0.1:6380> TS.RANGE temperature:3:11 1651708840000 1651708890000
1) 1) (integer) 1651708850000
   2) "27"
2) 1) (integer) 1651708860000
   2) "28"
3) 1) (integer) 1651708870000
   2) "29"
```


# License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FRedEpochDB%2FRedEpochDB.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FRedEpochDB%2FRedEpochDB?ref=badge_large)

# Disclaimers
When you use this software, you have agreed and stated that the author, maintainer and contributor of this software are not responsible for any risks, costs or problems you encounter. If you find a software defect or BUG, ​​please submit a patch to help improve it!

