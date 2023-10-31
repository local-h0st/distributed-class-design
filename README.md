# 分布式系统 课程设计实验

后面的实验专注于组播，所以我暂时把单播UDP通信部分给删了，最后一个具有单播功能的commit是[0ec5f3](https://github.com/local-h0st/distributed-class-design/commit/0ec5f349d5814dd00be373986598db1619cc03ce)。此外Windows多网卡环境下收不到组播消息，网上搜了一整个下午，问了ChatGPT还是没法解决，不知道是多网卡的问题还是golang的问题，无奈之下只能开linux来继续实验。

为在不同的程序间传递消息并且避免二义性，我们决定采用json格式来作为消息规范。一份作业的消息结构如下：

```go
type HOMEWORK_INFO struct {
	Name  string  // 姓名，例如"张三"
	ID    string  // 学号，例如"5712xxxx"
	Seq   string  // 作业序号，例如"1"
	Grade string  // 成绩，例如"96"
	Tag   string  // 是否按时提交，例如"Yes"或"No"
}
```

使用string类型是未来保证兼容性，因为不同语言的客户端解析json时，对不同数据类型可能有不同的表示方法，例如bool类型区分大小写（true和True）等，此外，考虑到作业序号可能包含非数字字符、成绩评定也可能使用ABCD而飞分数，故更加应该采用string类型。

我选择使用Golang构建一份消息：

```go
h := HOMEWORK_INFO{
	Name:  "张三",
	ID:    "57123456",
	Seq:   "1",
	Grade: "90",
	Tag:   "Yes",
}
s, _ := json.Marshal(h)
fmt.Println(string(s))
```

如果选择使用Python构建消息，那么一个可能的示例如下：

```python
import json

class HOMEWORK_INFO:
    def __init__(self, name, id, seq, grade, tag):
        self.Name = name
        self.ID = id
        self.Seq = seq
        self.Grade = grade
        self.Tag = tag

# 创建一个HOMEWORK_INFO对象
homework = HOMEWORK_INFO("张三", "57123456", "1", "90", "Yes")

# 将对象转换为JSON字符串
json_str = json.dumps(homework.__dict__)

# 打印JSON字符串
print(json_str)
```

示例json记录如下：

```json
{"Name": "张三", "ID": "5712xxxx", "Seq": "1", "Grade": "96", "Tag": "Yes"}
```

