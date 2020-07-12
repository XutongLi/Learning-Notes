# Lec10

EC2适用于web server，因为svr是无状态的。

EC2不利于应用于DB，因为DB与数据在一个硬盘上，如果硬盘故障，就会丢失数据

S3用于存储：定期建立快照

EBS：用于弹性存储，通过复制容错，使用的是链式复制（database mount EBS）

EC2不能用来shared，任何时刻，仅有一个EC2实例虚拟机可以挂载一个EBS volume

每个EBS volume只可以被一个EC2 instance使用

***

AZ就是一个数据中心

**transaction**： a way of wrapping multiple operations on maybe different pieces of data and declare in that entire sequence of operations should appear atomic to anyone else who is reading or writing the data

执行事务前会锁定对应数据，事务提交后，即持久化后，数据的锁释放

***

Aurora快的原因：网络中传输的仅仅是log record（之前是传送data pages）、quorum scheme

they really needed to have fast re-replicate them that is of one server seems permanently dead we'd like to be able to generate a new replica as fast as possible from the remaning replicas

想要写入N个副本，但只等待W个响应，读取只等待R个响应，W+R=N+1（不用等待慢的或故障的server）

***

读写分离：读仅从一个数据库执行

只读数据库会有一定延迟，但是可以接受