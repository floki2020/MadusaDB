<?php
    include "MadusaDB.class.php";

    $MadusaDB = new MadusaDB();
    //字符串类型操作

        //添加数据
        $string = $MadusaDB->key("MadusaDB")->add("MadusaDB");
        dump($string);
        //bool(true)

        //查找数据
        $string = $MadusaDB->key("MadusaDB")->find();
        dump($string);
        //string(8) "MadusaDB"

        //删除数据
        $status = $MadusaDB->key("MadusaDB")->del();
        dump($status);
        //bool(true)

    //数组操作(仅支持二位数组)

        $array['name']    =  "MadusaDB";
        $array['version'] =  "beta";
        //增加数组
        $status = $MadusaDB->key("array")->zadd($array);
        dump($status);
        //bool(true)

        //查找数组
        $array = $MadusaDB->key("array")->zfind();
        dump($array);
        /*  array(2) {
        *      ["name"] => string(8) "MadusaDB"
        *      ["version"] => string(3) "beta"
        *  }
        */

        //删除数组
        $status = $MadusaDB->key("array")->zdel();
        dump($status);
        //bool(true)

    //计数器操作

        //数字递增
        $int = $MadusaDB->incr("Madusa_num");
        dump($int);
        
        //数字递减
        $int = $MadusaDB->decr("Madusa_num");
        dump($int);
    //获取MadusaDB信息
        $info = $MadusaDB->info();
        dump($info);
        
        /*
        *   array(3) {
        *       ["keys"] => int(0)
        *       ["total_commands"] => int(10)
        *       ["version"] => string(3) "0.3"
        *   }
        */
 ?>
