# NAT Forward

Use for intranet penetration.

## how it works

```
+-----------------------------------+
|        +----------------+         |
|        | custom_client  |         |
|        +----+------+----+         |
| Internet    |      ^              |
|             |      |              |
|        +----v------+----+         |
|        | forward_server |         |
|        +----+------+----+         |
|             |      ^              |
+------------ | ---- | -------------+
              |      |
+------------ | ---- | -------------+
|        +----v------+----+         |
|        | forward_client |         |
|        +----+------+----+         |
| Behind      |      ^              |
| Nat         |      |              |
|        +----v-----------+         |
|        | actual_ser|er  |         |
|        +----------------+         |
+-----------------------------------+

```

## ref

[golang手把手实现tcp内网穿透代理(3)快速编码测试](https://www.jianshu.com/p/e79fe205f3e0)