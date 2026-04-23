# hash

## 密码存储hash

命名：

- `pwdSHA256`: `sha256('password text')`，服务端接收该类型的参数
- `key`: `argon2('pwdSHA256', 'salt')`，算法派生的密钥
    - argon2算法可以根据算法参数、密码和盐派生密钥（相同环境下，重复派生的结果相同）
    - `salt`: 程序随机生成
- `pwdArgon2`: `encode(key)`
    - 编码后结构示意: `argon2id.v=19,m=65536,t=3,p=1.[saltHex].[keyHex]`

密码存储不使用明文，前端传参数`pwdSHA256`、服务端保存的则是`pwdArgon2`

前端执行hash的目的是即使恶意攻击窃取了http请求参数，也无法恢复出密码明文，进而攻击者无法使用我们的ui；
后端执行hash的目的是即使恶意攻击窃取了数据库数据，也无法恢复出密码hash，进而无法调用接口。

todo：介绍本节与access token使用不同hash算法的原因和考量
