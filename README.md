# Ant
网络爬虫服务

这个是源码
需要这个服务的可以调用https://cninct.com/Ant
以POST方式提交需要爬取的信息

示例：

{

    "url": "http://company.zhaopin.com/zhengzhou/210500/",//需要爬取的目标网址

    "ant": [//需要爬取的具体信息，可以同时爬取多个同页面的目标，可为空

        {

            "parend": {//目标节点的父节点属性要求，可为空

                "class": "fleft checkjobs width280"

            },

            "grandpa": {//目标节点的父节点的父节点的属性要求，可为空

                "class": "jobs-list-box"

            },

            "curr": {},//目标节点的属性要求，可为空

            "currAtom": "A",//目标节点标签属性要求，可为空

            "attr": [//爬取信息默认爬取text属性，如果有其它属性也需要爬取可以在这里添加

                "href"

            ],

            "utf8": true //爬取信息的编码格式，默认为utf8属性，如果需要gbk请设置为false

        }

    ]

}