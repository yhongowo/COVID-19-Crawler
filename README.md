# 2019-nCov-Crawler
COVID-19/2019-nCoV 新型冠状病毒实时爬虫，采用Go，gocolly，goquery。数据来源：[丁香园](https://ncov.dxy.cn/ncovh5/view/pneumonia "丁香园")

## 介绍
采用MongoDB作为数据库，使用了时序数据库。每隔一个小时进行爬虫，将数据处理后持久化，在tmp目录里存放了上次爬虫的数据。

## 部署
- 环境：Golang1.18
- git clone到本地
- /service/db.go中配置数据库
```
const URI = "mongodb://localhost:27017"
const DATABASE = "2019-nCov"
```
- $run main.go 启动服务
- 也可以在release中下载exe文件试运行
- 