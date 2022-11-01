# Jenv
Go实现的简易Windows多版本Java环境变量管理小工具
需要手动在Path中添加`%JAVA_HOME%\bin`

功能：
jenv list               ：列出已添加的Java目录
jenv add [tag] [path]   ：添加Java目录的记录
jenv use [tag]          ：切换到tag对应的Java版本
jenv del [tag]          ：删除tag对应记录

例：
jenv add jdk8 C:\env\jdk-8    添加记录
jenv use jdk8                 切换到jdk8
