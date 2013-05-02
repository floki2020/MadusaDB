<?php
    include "xuncache.class.php";

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
 ?>
