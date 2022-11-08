# 海量用户即时通讯系统


## 需求分析

- 用户注册
- 用户登录
- 显示在线用户列表
- 群聊(广播)
- 点对点聊天
- 离线留言



## 功能实现

### 显示客户登录菜单

功能:能够正确的显示客户端的菜单

### 完成用户登录

#### 客户端

- 接收输入的uname和upwd

- 发送uname和upwd

- 接收服务端返回的结果，判断成功还是失败，显示相对应的页面

- 关键的问题是，**怎样组织发送的数据**

  - 设计消息协议

    - **`Message Struct`**
      - `type string`
      - `data string`
    - **`LoginMes struct`**
      - `uname string`
      - `upwd string`
    - 消息类型**-->`type string`**
    - 消息结构体序列化后塞给**`Message.Data`**

  - **发送的流程**

    1. 先创建一个Message结构体
    2. 设置一个消息类型 例如 `Mes.Type`=登录消息类型
    3. `mes.Data`=登录消息的内容(序列化之后)
    4. 对Message结构体再进行序列化
    5. 在网络传输中，最担心的就是丢包问题，**解决**：
       - 先给服务器发送message的长度（有多少个字节）
       - 再发送消息本身

    

  - **接收的流程**

    1. 接收到客户端发送的长度
    2. 根据接收到的长度，再接收消息本身
    3. 接收时，要判断实际接收到的消息内容是否等于`len`
    4. 如果不相等，就有纠错协议
    5. 取到后可以反序列化成`Message Struct`类型
    6. 取出`message.Data(string)`>`loginmessage`
    7. 取出`loginmessage`的`uname`和`upwd`
    8. 这时就可以比较
    9. 根据比较结果，返回mess
    10. 发送给客户端



### 代码步骤分析

1. 完成客户端可以发送消息长度，服务器端可以正常接收到该长度

   - 先确定Message的格式和结构 

2.  完成客户端可以发送消息本身，服务器端可以正常接收消息，并根据客户端发送的消息(`LoginMes`)，判断用户的合法性，并返回相应的`LoginResMes`

   思路分析：

   - 让**[客户端]**发送消息本身
   - **[服务器端]**收到消息，并反序列化成对应的消息结构体
   - **[服务器端]**根据反序列化后的消息，判断登录的用户是否为合法用户，返回`LoginResMes`
   - **[客户端]**解析返回的`LoginResMes`，显示对应界面
   - 这里我们需要做一些函数的封装



## 服务器代码改进

**步骤**

- #### 现将分析出的文件创建，放入相应的包中

  ![](F:\GoProject\pro_connection\pics\server.png)

- #### 项目文件结构

  ```s
  │  go.mod
  │  readme.md
  │  新建 XLS 工作表.xls
  │
  ├─.vscode
  │      launch.json
  │      settings.json
  │
  ├─client
  │      login.go
  │      main.go
  │      utils.go
  │
  ├─common
  │  └─message
  │          message.go
  │
  └─server
      ├─main
      │      main.go
      │      processor.go
      │
      ├─model
      ├─process
      │      smsProcess.go
      │      userProcess.go
      │
      └─utils
              utils.go
  ```

- 将相应功能实现放入指定文件中



## 客户端代码改进

- #### **客户端程序结构**

  ![](F:\GoProject\pro_connection\pics\client.png)

- #### 项目文件结构

  ```s
  ├─main
  │      main.go
  │
  ├─model
  ├─process
  │      server.go
  │      smsProcess.go
  │      userProcess.go
  │
  └─utils
          utils.go
  ```





## 实现功能：完成用户登录

- 在`Redis`手动添加测试用户（后面通过程序注册用户）

  ```
  127.0.0.1:6379> hset users 1 "{\"userid\":1,\"upwd\":\"199866\",\"uname\":\"Hud\"}"
  (integer) 0
  127.0.0.1:6379> hset users 2 "{\"userid\":2,\"upwd\":\"199866\",\"uname\":\"Monica\"}"
  (integer) 1
  ```

  如果输入的用户名密码正确在`Redis`中存在则登录，否则退出系统 :-)

- 给出相应的提示信息：
  - 用户不存在或密码错误
  - 你也可以重新注册，再登录



## 实现功能：完成用户的注册

- ##### 完成用户注册功能，将用户信息录入到`Redis`中

  - 在message中定义两个新的消息类型
    - `struct RegisterMes`
    - `struct RegisterResMes`
  - 在客户端接收用户的输入(界面)
  - 在客户端的`UserProcess.go`中编写一个Register方法，完成请求注册的功能
  - 在服务器端的`Server/model/userDao.go`中新增一个方法，在`redis`中实现数据的添加



## 实现功能：完成登录时能返回当前在线的用户

1. 在服务器端维护一个**`onlineUsers`**，结构是**`map [int]*UserProcess`**
2. 创建一个新的文件：**`userMgr.go`** ，完成功能：对**`usermap`**的增删改查
3. 在**`LoginResMess`**结构体中增加一个字段：**`Users []int`** ,将在线用户的id返回

- #### **思路**

  1. 当一个用户A上线，服务器就把A用户的上线信息，推送给所有的在线的用户
  2. 客户端也需要维护一个`map`，map中记录了他们的好友(目前就是所有人) `map[int]User`
  3. 客户端和服务器的通讯通道，要依赖`serverProcessMes`协程





## 实现功能：用户进行消息群发

1. 新增一个消息结构体`SmsMes`...

2. 新增一个`model CurUser`

   ```go
   type CurUser struct {
   	Conn net.Conn
   	message.User
   }
   ```

3. 在`SmsProcess.go`增加相应的方法`BroadCast`，发送一个群聊消息

4. 服务器接收到消息，推送给所有的在线用户(发送者除外)。

5. 在服务器端接收到`SmsMes`

6. 在服务器`SMSProcess.go`文件中增加一个群发消息的方法

7. 在客户端还要增加去处理服务器端转发的群发消息
