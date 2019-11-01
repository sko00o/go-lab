# ISSUES

Compare the differences on IX (intention lock).

- SELECT ... IN ...

```sql
-- tx1
START TRANSACTION;
SELECT * FROM seats WHERE seat_no IN (3, 4) FOR UPDATE;
```

```text
+---------+--------+
| seat_no | booked |
+---------+--------+
|       3 | NO     |
|       4 | NO     |
+---------+--------+
```

```sql
-- tx2
SELECT object_name, index_name, lock_type, lock_mode, lock_data
FROM performance_schema.data_locks WHERE object_name = 'seats';
```

```text
+-------------+------------+-----------+---------------+-----------+
| object_name | index_name | lock_type | lock_mode     | lock_data |
+-------------+------------+-----------+---------------+-----------+
| seats       | NULL       | TABLE     | IX            | NULL      |
| seats       | PRIMARY    | RECORD    | X,REC_NOT_GAP | 3         |
| seats       | PRIMARY    | RECORD    | X,REC_NOT_GAP | 4         |
+-------------+------------+-----------+---------------+-----------+
```

- SELECT ... BETWEEN ... AND ...

```sql
-- tx1
START TRANSACTION;
SELECT * FROM seats WHERE seat_no BETWEEN 3 AND 4
```

```text
+---------+--------+
| seat_no | booked |
+---------+--------+
|       3 | NO     |
|       4 | NO     |
+---------+--------+
```

```sql
--tx2
SELECT object_name, index_name, lock_type, lock_mode, lock_data
FROM performance_schema.data_locks WHERE object_name = 'seats';
```

```text
+-------------+------------+-----------+---------------+-----------+
| object_name | index_name | lock_type | lock_mode     | lock_data |
+-------------+------------+-----------+---------------+-----------+
| seats       | NULL       | TABLE     | IX            | NULL      |
| seats       | PRIMARY    | RECORD    | X,REC_NOT_GAP | 3         |
| seats       | PRIMARY    | RECORD    | X             | 4         |
| seats       | PRIMARY    | RECORD    | X             | 5         |
+-------------+------------+-----------+---------------+-----------+
```
