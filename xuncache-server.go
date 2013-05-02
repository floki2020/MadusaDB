package main

import (
	handle "./handle"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var Store_int = make(map[string]uint64)                          //计数器
var Store_str = make(map[string]string)                          //字符串
var Store_map = make(map[string]map[string]interface{})          //map(array)
var Store_list = make(map[string][]map[string]interface{})       //列表
var Store_list_chain = make(map[string][]map[string]interface{}) //列表
var new_config = make(map[string]string)                         //配置文件
var total_commands_processed uint64                              //总处理数
var err error

const (
	IP          string = "127.0.0.1"
	PORT        string = "1200"
	CONFIG_FILE string = "xuncache.conf"
	VERSION     string = "0.4"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始内容
	fmt.Printf("Server started, xuncache version %s\n", VERSION)
	//读取配置文件
	new_config = read_config()
	//创建服务端
	tcpAddr, err := net.ResolveTCPAddr("tcp4", new_config["bind"]+":"+new_config["port"])
	fmt.Printf("The server is now ready to accept connections on %s:%s\n", new_config["bind"], new_config["port"])
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	//输出状态
	go func() {
		for {
			//自身占用
			now_time := time.Now().Format("2006-01-02 15:04:05")
			map_num := len(Store_map) + len(Store_str) + len(Store_int) + len(Store_list)
			fmt.Printf("[%s]DB keys: %d ,total_commands: %d ,concurrent :%d \n", now_time, map_num, total_commands_processed, runtime.NumGoroutine()-5)
			time.Sleep(2 * time.Second)
		}
	}()
	//数据处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

//处理数据
func handleClient(conn net.Conn) {
	var data_str, lfield string
	var limit_start, limit_stop, index uint64
	var limit_bool bool
	//标记结束连接
	defer conn.Close()
	defer fmt.Print("Client closed connection\n")
	ipAddr := conn.RemoteAddr()
	fmt.Printf("Accepted %s\n", ipAddr)
	for {
		//清空数据
		data_str, lfield = "", ""
		limit_start, limit_stop, index = 0, 0, 0
		var data_map = make(map[string]interface{})
		var back = make(map[string]interface{})
		//获取数据
		var buf [1024]byte
		n, _ := conn.Read(buf[0:])
		b := []byte(buf[0:n])
		if len(b) < 1 {
			return
		}
		total_commands_processed++ //记录处理次数
		js, err := handle.NewJson(b)
		if err != nil {
			back["error"] = true
			back["point"] = err.Error()
			rewrite(back, conn)
			return
		}
		pass, _ := js.Get("Pass").String()
		if pass != new_config["password"] && len(new_config["password"]) > 1 {
			fmt.Printf("Encountered a connection password is incorrect Accepted %s\n", ipAddr)
			back["error"] = true
			back["point"] = "password error!"
			rewrite(back, conn)
			return
		}
		//获取协议
		protocol, _ := js.Get("Protocol").String()
		//数据处理
		if protocol == `set` || protocol == `zset` || protocol == `push` || protocol == `lget` || protocol == `lincr` || protocol == `lsondelete` {
			switch protocol {
			case `set`:
				data_str, _ = js.Get("Data").String()
			case `zset`:
				data_map, _ = js.Get("Data").Map()
			case `push`:
				data_map, _ = js.Get("Data").Map()
			case `lget`:
				limit_start = js.Get("Start").Uint64()
				limit_stop = js.Get("Stop").Uint64()
				limit_bool = js.Get("Limit").Bool()
			case `lincr`:
				lfield, _ = js.Get("Field").String()
				index = js.Get("Index").Uint64()
			case `ldecr`:
				lfield, _ = js.Get("Field").String()
				index = js.Get("Index").Uint64()
			case `lsondelete`:
				index = js.Get("Index").Uint64()
			default:
				fmt.Print("error type \n")
				break
			}
			//数据判断
			if (data_map == nil && len(data_str) < 1) && (protocol == `set` || protocol == `zset` || protocol == `push`) {
				fmt.Print("There is no data \n")
				break
			}

		}
		//获取key
		key, _ := js.Get("Key").String()
		if len(key) < 1 && protocol != "info" {
			fmt.Printf("Error agreement is key %s\n", key)
			back["error"] = true
			back["point"] = "Please input Key!"
			rewrite(back, conn)
			break
		}

		//协议判断 处理
		switch protocol {
		case `set`:
			//字符串设置
			Store_str[key] = data_str
			back["status"] = true
			break
		case `zset`:
			//数组设置
			Store_map[key] = data_map
			back["status"] = true
			break
		case `push`:
			//list map添加(目前没加锁...目前测试性能没碰到过需要原子锁)
			data_map["_id"] = len(Store_list[key]) + 1
			Store_list[key] = append(Store_list[key], data_map)
			Store_list_chain[key] = Store_list[key]
			back["_id"] = data_map["_id"]
			back["status"] = true
			break
		case `get`:
			//获取字符串
			back["data"] = Store_str[key]
			back["status"] = true
			break
		case `zget`:
			//获取数组
			back["data"] = Store_map[key]
			back["status"] = true
			break
		case `lget`:
			//这里是list map 查询
			if limit_start > 0 {
				limit_start--
			}
			Store_list_nums := uint64(len(Store_list[key]))
			if limit_start > Store_list_nums {
				back["status"] = false
				break
			}
			if limit_start+limit_stop > Store_list_nums {
				limit_stop = Store_list_nums - limit_start
			}
			//分页查询
			if limit_stop > 0 {
				back["data"] = Store_list[key][limit_start : limit_start+limit_stop]
			} else if limit_start > 0 && limit_stop == 0 && limit_bool == false {
				//分页查询--容错
				back["data"] = Store_list[key][Store_list_nums-limit_start : limit_start+limit_stop]
			} else if limit_bool == true {
				//查询单条
				back["data"] = Store_list[key][limit_start:]
			} else {
				//查询所有
				back["data"] = Store_list[key][0:]
			}
			back["status"] = true
			break
		case `delete`:
			//删除字符串
			delete(Store_str, key)
			back["status"] = true
			break
		case `zdelete`:
			//删除数组
			delete(Store_map, key)
			back["status"] = true
			break
		case `ldelete`:
			//删除list map
			delete(Store_list, key)
			back["status"] = true
			break
		case `lsondelete`:
			//删除list map(此处用销毁内存,并移动元素实现.不完美,后期考虑用双向链表实现)
			if int(index) > len(Store_list[key]) {
				back["status"] = false
				break
			}

			Store_list[key][index-1] = nil
			Store_list[key] = append(Store_list[key][:int(index)-1], Store_list[key][int(index):]...)
			back["status"] = true
			break
		case `lincr`:
			//list 单字段计数器
			if int(index) > len(Store_list[key]) {
				back["status"] = false
				break
			}
			if nums, ok := Store_list[key][int(index)-1][lfield].(float64); ok {
				back["status"], Store_list[key][int(index)-1][lfield] = nums+1, nums+1
			}
			break
		case `ldecr`:
			//list 单字段计数器
			if int(index) > len(Store_list[key])-1 {
				back["status"] = false
				break
			}
			if nums, ok := Store_list[key][int(index)][lfield].(float64); ok {
				if nums < 1 {
					back["status"], Store_list[key][int(index)][lfield] = 0, 0
				} else {
					back["status"], Store_list[key][int(index)][lfield] = nums-1, nums-1
				}
			}
			break
		case `incr`:
			//计数器++
			Store_int[key]++
			back["data"] = Store_int[key]
			back["status"] = true
			break
		case `decr`:
			//计数器--
			if Store_int[key] == 0 {
				back["data"] = 0
				back["status"] = true
				break
			}
			Store_int[key]--
			back["data"] = Store_int[key]
			back["status"] = true
			break
		case `info`:
			var info = make(map[string]interface{})
			info["version"] = VERSION
			info["keys"] = len(Store_map) + len(Store_str) + len(Store_int) + len(Store_list)
			info["total_commands"] = total_commands_processed
			back["data"] = info
			back["status"] = true
			break
		default:
			back["status"] = false
			fmt.Print("error protocol \n")
			break
		}
		//返回内容
		rewrite(back, conn)
	}
}

//读取配置文件
func read_config() (new_config map[string]string) {
	var config = make(map[string]string)
	dir, _ := path.Split(os.Args[0])
	os.Chdir(dir)
	path, _ := os.Getwd()
	config_file, err := os.Open(path + "/" + CONFIG_FILE) //打开文件
	defer config_file.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Print("Can not read configuration file. now exit\n")
		os.Exit(0)
	}
	buff := bufio.NewReader(config_file) //读入缓存
	//读取配置文件
	for {
		line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil {
			break
		}
		rs := []rune(line)
		if string(rs[0:1]) == `#` || len(line) < 3 {
			continue
		}
		str_type := string(rs[0:strings.Index(line, " ")])
		detail := string(rs[strings.Index(line, " ")+1 : len(rs)-1])
		config[str_type] = detail
	}
	//再次过滤 (防止没有配置文件)
	return verify(config)
}

//写入数据
func rewrite(back map[string]interface{}, conn net.Conn) {
	jsback, _ := json.Marshal(back)
	//返回内容
	conn.Write(jsback)
}

//验证配置文件
func verify(config map[string]string) (config_bak map[string]string) {
	if len(config["bind"]) < 3 {
		config["bind"] = IP
	}
	if len(config["port"]) < 1 {
		config["port"] = PORT
	}
	return config
}

func item_expired() {

}

//输出错误信息
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
