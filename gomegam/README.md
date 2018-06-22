# GoMegam-IoT


As the name depicts, the project revolves around  Internet of Things(IoT), an umbrella term, which is basically bunch of devices communicating and storing data through the internet.  Every device generates some sort of data all the time, these data if well filtered, can be used for further analysis to get valuable insights.  

![Concept](https://github.com/gophergala/gomegam/blob/master/conf/gomegam.jpg)

## How can these data be distributed, fail-proof and redundant? 

Say hi to blockchain. We use blockchain, the underlying technology that runs bitcoins to do the job. To put things short, block chain is a technology which is decentralized and distributed, everyone in the network contains all the information about the data when ever there is a new transaction that takes place in the network. This way, we eliminate decentralized system and increase redundancy of data. WIN-WIN.

## Target user?

It can be any one who would want to use the blockchain technology and be distributed. A healthcare firm? Yes, they can inturn connect their devices to the blockchain network and hence not worry about losing data. hey remember, is it still encrypted!  an energy firm? Absolutely, Send filtered data to the network, analysed later would give important insights. 

## Why do we use go?

We use go because of two simple reasons

1. It is awesome

2. It really is awesome. 

on a serious note, Go can be used effectively because of its concurrency and low memory usage. When there are way too many devices trying to talk to each other, or to a main network, concurrency comes in handy. Writing code in go is a breeze!

##GomegamIot

A blockchain aware daemon built on beego, and go channels


```
git clone https://github.com/gophergala/gomegam

make test

```

```
$ iot
2015/01/26 03:03:19 Conf: /home/ram/code/gomegam/conf/iot.conf
2015/01/26 03:03:19 [config.go:284] [W] open /home/ram/code/gomegam/conf/app.conf: no such file or directory 
iot version 0.1.0.

Usage: iot command [args]

Available commands:
  help
  start
  version

Use iot help <commandname> to get more information about a command.

$ iot start

```

# License

|                      |                                          |
|:---------------------|:-----------------------------------------|
| **Author:**          | Rajthilak (<rajthilak@megam.co.in>)
|                      | Thomas Alrin (<thomasalrin@megam.co.in>)
|		       	       | KishorekumarNeelamegam (<nkishore@megam.co.in>)
|		       	       | Yeshwanth Kumar (<getyesh@megam.co.in>)
| **Copyright:**       | Copyright (c) 2013-2015 Megam Systems.
| **License:**         | Apache License, Version 2.0

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 


