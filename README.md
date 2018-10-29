# dgraph-bench

A benchmark program for dgraph.

## Benchmark Results

All tests are on servers with multi-way Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz,
32 virtual cores on each server.

Each server has 64G memory, and a 500GB SATA SSD:

| Parameter            | Evaluation           |
|----------------------|----------------------|
| Max Sequential Read  | 560 MBps             |
| Max Sequential Write | 530 MBps             |
| 4KB Random Read      | 95,000 IOPS          |
| 4KB Random Write     | 84,000 IOPS          |

### Simple Node Insertion

Af the first step, we tried to insert single node into dgraph, in different
cluster configurations.

#### 1 zero, 1 alpha, 100 concurrent goroutines, no index

![image](1zero-1alpha-simple-insert.png)

The performance is not bad, considered that the concurrency is 100 and there is
a 20ms cool down for each call. Meanwhile the CPU load is low:

![image](1zero-1alpha-simple-insert-cpu.png)

#### 1 zero, 3 alpha, 100 concurrent goroutines, no index

![image](1zero-3alpha-simple-insert.png)

Now we see a poor performance. The CPU load is still not high:

Machine 1 (with zero and ratel):

![image](1zero-3alpha-simple-insert-cpu-1.png)

Mchine 2 and 3:

![image](1zero-3alpha-simple-insert-cpu-2.png)
![image](1zero-3alpha-simple-insert-cpu-3.png)

After one hour of running, we noticed a performance jump, followed by the same
performance drop.

![image](1zero-3alpha-simple-insert-longtime.png)
