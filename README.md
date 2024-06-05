# Mages引擎工具集
从头开始写的Go版，算是对以前的工具进行一次整合和完善，进度会很慢，结构会经常大改  
目标是完成尽可能通用的Mages工具集，目前仅以脚本为主

## 适配游戏
- 理论支持所有Mages引擎的游戏
- 所有的MES(msb)、SC3(scx)脚本都可正常导出导入

## Usage
```
  -charset string
        [script.optional] Character set containing only text. Must be utf8 encoding. Choose between "charset" and "tbl"
  -debug int
        [optional] Debug level
            0: Disable debug mode
            1: Show info message
            2: Show warning message (For example, the character table is missing characters)
            3: Not implemented
  -export
        [optional] Export mode. Support folder export
  -format string
        [script.required] Format of script export and import. Case insensitive
            NPCSManager format: "Npcs"
            NPCSManager Plus format: "NpcsP" (default "Npcs")
  -import
        [optional] Import mode
  -input string
        [optional] Usually the import mode requires
  -output string
        [required] Output file or folder
  -skip
        [script.optional] Skip repeated characters in the character table. (default true)
  -source string
        [required] Source files or folder
  -tbl string
        [script.optional] Text in TBL format. Must be utf8 encoding. Choose between "charset" and "tbl"
  -type string
        [required] Source file type.
            Mages Script: "script"
                Supported MES(msb), SC3(scx)
            Diff Binary File: "diff"
                Diff input and output file


```
### Example

```shell
# 导出文件夹所有，使用tbl码表，格式为NpcsP，不跳过码表中相同字符，开启debug模式为2
MagesTools -type=script -export -skip=false -debug=2\
  -format=NpcsP \
  -tbl=./data/CC/MJPN.txt \
  -source=./data/CC/script/mes00 \
  -output=./data/CC/txt


# 导出文本，使用tbl码表，格式为NpcsP，跳过码表中相同字符
MagesTools -type=script -export -skip=true \
  -format=NpcsP \
  -tbl=./data/CC/MJPN.txt \
  -source=./data/temp/1.msb \
  -output=./data/temp/1.msb.txt 

  
# 导入文本，使用tbl码表，格式为NpcsP，跳过码表中相同字符
MagesTools -type=script -import -skip=false \
  -format=NpcsP \
  -tbl=./data/CC/MJPN.txt \
  -source=./data/temp/1.msb \
  -input=./data/temp/1.msb.txt \
  -output=./data/temp/1.msb.txt.msb

# RNE使用以下参数
# 导出文本，使用charset码表，格式为Npcs，不跳过码表中相同字符
MagesTools -type=script -export -skip=false \
  -format=Npcs \
  -charset=./data/RNE/Charset_PSV_JP.utf8 \
  -source=./data/temp/1.msb \
  -output=./data/temp/1.msb.txt 

  
# 导入文本，使用charset码表，格式为Npcs，不跳过码表中相同字符
MagesTools -type=script -import -skip=false \
  -format=Npcs \
  -charset=./data/RNE/Charset_PSV_JP.utf8 \
  -source=./data/temp/1.msb \
  -input=./data/temp/1.msb.txt \
  -output=./data/temp/1.msb.txt.msb
  
# 对比文件
MagesTools -type=diff \
  -input=./data/temp/1.msb \
  -output=./data/temp/1.msb.txt.msb
```

## Script
### 格式
目前的格式为NPCSManager的优化版
- 删除`name`后的`[1x01][1x02]`，仅用`:[`和`]:`标记name
- 删除`]:`后的半角空格
- 保留字节数据均采用`0x`开头，如`[0x04A01414]`
- 删除`color`的特殊标记`<#`和`#>`，仅用字节标记，如`[0x04A01414][0x00]`
- 增加对`EvaluateExpression`表达式的简单字节解析，如`[0x15290AA4B51414008100][0x00]`，可能存在未知bug  
```
[0x0F][0x1100CC][0x04A01414][0x00]『白い光が見えた』[0x15290AA4B51414008100][0x00][0x03][0xFF]
[0x0F][0x110026][0x04A01414][0x00]『耳鳴りのような音が聞こえた』[0x15290AA4B51414008100][0x00][0x08][0xFF]
[0x0F][0x1100F2]勘違いだと笑ってしまうにはあまりに多くの者たちが[0x1F]体験してしまったこの現象は、原因不明のまま語り継がれ、[0x1F]地震のおかしさを疑う者の手助けをすることとなった。[0x15290AA4B51414008100][0x00][0x08][0xFF]
[0x0F][0x110118]そして、噂にまみれた地震から６年経った２０１５年。[0x15290AA4B51414008100][0x00][0x08][0xFF]
[0x0F][0x1100F2]新しく生まれ変わりつつある渋谷の街で、[0x1F]地震とは別の事件が世間の注目を集めようとしていた。[0x15290AA4B51414008100][0x00][0x08][0xFF]
[0x0F][0x110118]２０１５年９月７日（日）夜[0x15290AA4B51414008100][0x00][0x08][0xFF]
:[男性]:「はい、ではいつも通り３分くらい募集をかけるんで、適当によろです」[0x03][0xFF]
そう言った途端に、コメント欄が一気に流れ出した。[0x03][0xFF]
流れ具合を数秒間確認していると、『ハルちゃんの熱愛報道はいつ？』との依頼を見つけ、[0x09]大谷[0x0A]おおたに[0x0B][0x09]悠馬[0x0A]ゆうま[0x0B]は思わず微笑んだ。[0x03][0xFF]
狙い通りだ。[0x03][0xFF]
依頼の大半はイケメン俳優か女性アイドルに関することだから、視聴者の傾向を読むことは馬鹿みたいに簡単だった。[0x03][0xFF]
問題は、どの人物の名前が挙がるかということで、こればかりは運とその人物の人気による。[0x03][0xFF]
が、[0x04280AA0][0x2D14][0x00]ハルちゃん[0x04800000][0x8113][0x8113]確かなんとかハルコとかいったか[0x8113][0x8113]ならば大丈夫だ。[0x03][0xFF]
先日、行きたくもないイベントに行って、[0x09][0x1E]直接見て来た[0x0A][0x8117][0x8117][0x8117][0x8117][0x8117][0x8117][0x0B]ばかりだ。[0x03][0xFF]
:[大谷]:「……よし」[0x03][0xFF]
```

## 计划
- 支持更多格式

## 更新日志

### 2024.6.5
- 修复表达式结束判断
- 复无文本的脚本解析问题
- 修复windows下导出文件夹无写入问题 
- (以上问题由 [Fluchw](https://github.com/wetor/MagesTools/issues/5) 发现)

### 2022.10.21
- 修复':'字符后为字节数据（`:[0xFF]`）导致Encode错误的问题 ([kurikomoe](https://github.com/kurikomoe)发现)

### 2022.3.21
- 支持SC3(scx)脚本的文本导入导出
- 支持文件夹导出（暂不支持导入）
- 增加简单的log
- 优化代码细节

### 2022.3.20 2
- 重构代码结构以支持更多导出格式
- 支持NPCSManager格式的导出与导入
  - 支持NPCSManager导出文本的导入
  - NPCSManager无法导入此程序导出的文本，存在细微差异
- 支持命令行调用
- 增加帮助文档

### 2022.3.20
- 完成MES(msb)文本导入（简单实现）
- 调整导出格式

### 2022.3.19
- 基本框架设计
- 完成MES(msb)的文本导出
### 2022.3.18
- 开始


## 参考依赖
- [marcussacana](https://github.com/marcussacana) 的 [NPCSManager](https://github.com/marcussacana/NPCSManager)  
- [liaowm5](https://github.com/SteiensGate) 的 msb_tool.py
- [CommitteeOfZero](https://github.com/CommitteeOfZero) 的 [sc3ntist](https://github.com/CommitteeOfZero/sc3ntist) 以及 [SciAdv.Net](https://github.com/CommitteeOfZero/SciAdv.Net)