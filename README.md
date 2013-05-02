## xuncache
========
xuncache 是免费开源的NOSQL(内存数据库) 采用golang开发,简单易用而且 功能强大(就算新手也完全胜任)、性能卓越能轻松处理海量数据,可用于缓存系统.

目前版本 version 0.3

前期它是活跃的 更新很迅速

version 1.0版本前 作者不推荐用于生产环境

采用json协议 socket通信 --后期打算用bson

## 目前功能
========
-增加or设置(字符串和数组)

-查找数据(字符串和数组)

-删除数据(字符串和数组)

-计数器功能

-暂不支持key过期操作

支持 php 客户端 
## php代码示例
========

	$xuncache = new xuncache();
    //字符串类型操作

        //添加数据
        $string = $xuncache->key("xuncache")->add("xuncache");
        dump($string);
        //bool(true)

        //查找数据
        $string = $xuncache->key("xuncache")->find();
        dump($string);
        //string(8) "xuncache"

        //删除数据
        $status = $xuncache->key("xuncache")->del();
        dump($status);
        //bool(true)

    //数组操作(仅支持二位数组)

        $array['name']    =  "xuncache";
        $array['version'] =  "beta";
        //增加数组
        $status = $xuncache->key("array")->zadd($array);
        dump($status);
        //bool(true)

        //查找数组
        $array = $xuncache->key("array")->zfind();
        dump($array);
        /*  array(2) {
        *      ["name"] => string(8) "xuncache"
        *      ["version"] => string(3) "beta"
        *  }
        */

        //删除数组
        $status = $xuncache->key("array")->zdel();
        dump($status);
        //bool(true)

    //计数器操作

        //数字递增
        $int = $xuncache->incr("xuncache_num");
        dump($int);
        
        //数字递减
        $int = $xuncache->decr("xuncache_num");
        dump($int);
    //获取xuncache信息
        $info = $xuncache->info();
        dump($info);
        
        /*
        *   array(3) {
        *       ["keys"] => int(0)
        *       ["total_commands"] => int(10)
        *       ["version"] => string(3) "0.3"
        *   }
        */
	
## 性能测试(仅代表本机速度)
###不是专业测试 如果你有更好的测试结果欢迎提交

![](images/property.png?raw=true)

## 关于
- by [孙彦欣](http://weibo.com/sun8911879)
-    [更新日志](https://github.com/sun8911879/xuncache/blob/master/UPDATE.md)
- LICENSE: under the [BSD](https://github.com/sun8911879/xuncache/blob/master/LICENSE-BSD.md) License