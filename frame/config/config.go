package config

import(
	"flag"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"encoding/json"
)

// See handleFlags methods about these vars and their uses.
type Config struct{
	AbsPath		string
	ConfFn		string
	DBAdmMode	bool
	DBUser		string
	DBPass		string
	DBAddr		string
	MsgAddr		string
	Debug		bool
	DBName		string
	Addr		string
	PortNum		string
	CacheOpt	bool
	ServeFiles	bool
	Secret		string
}

var cli = Config{}

func New() *Config {
	args()
	ret := cli
	return &ret
}

func (c *Config) LoadFromFile() {
	cf, err := ioutil.ReadFile(filepath.Join(c.AbsPath, c.ConfFn))
	if err != nil {
		return
	}
	var conf_i interface{}
	err = json.Unmarshal(cf, &conf_i)
	if err != nil || conf_i == nil {
		return
	}
	conf, ok := conf_i.(map[string]interface{})
	if !ok {
		return
	}
	if db_adm_mode, ok := conf["db_admin_mode"].(bool); ok {
		c.DBAdmMode = db_adm_mode
	}
	fmt.Println("Overwriting some command line args with ones provided in the config file.")
	// Doh...
	if db_user, ok := conf["db_user"].(string); ok {
		c.DBUser = db_user
	}
	if db_pass, ok := conf["db_pass"].(string); ok {
		c.DBPass = db_pass
	}
	if db_addr, ok := conf["db_addr"].(string); ok {
		c.DBAddr = db_addr
	}
	if debug, ok := conf["debug"].(bool); ok {
		c.Debug = debug
	}
	if db_name, ok := conf["db_name"].(string); ok {
		c.DBName = db_name
	}
	if addr, ok := conf["addr"].(string); ok {
		c.Addr = addr
	}
	if port_num, ok := conf["port_num"].(string); ok {
		c.PortNum = port_num
	}
	if cache_opt, ok := conf["cache_opt"].(bool); ok {
		c.CacheOpt = cache_opt
	}
	if serve_files, ok := conf["serve_files"].(bool); ok {
		c.ServeFiles = serve_files
	}
	if secret, ok := conf["secret"].(string); ok {
		c.Secret = secret
	}
}

func args() {
	flag.StringVar(	&cli.AbsPath, 		"abs_path", 	"c:/gowork/src/github.com/opesun/nocrud", "absolute path")
	flag.StringVar(	&cli.ConfFn, 		"conf_fn", 		"config.json", 		"config filename")
	flag.BoolVar(	&cli.DBAdmMode, 	"db_adm_mode", 	false, 				"connect to database as an admin")
	flag.StringVar(	&cli.DBUser, 		"db_user", 		"", 				"database username")
	flag.StringVar(	&cli.DBPass, 		"db_pass", 		"", 				"database password")
	flag.StringVar(	&cli.DBAddr, 		"db_addr", 		"127.0.0.1:27017", 	"database address")
	flag.StringVar(	&cli.MsgAddr, 		"msg_addr", 	"", 				"message passing server address")
	flag.StringVar(	&cli.DBName, 		"db_name", 		"nocrud", 			"db name to connect to")
	flag.BoolVar(	&cli.Debug, 		"debug", 		true, 				"debug mode")
	flag.StringVar(	&cli.PortNum, 		"p", 			"80", 				"port to listen on")
	flag.StringVar(	&cli.Addr, 			"addr", 		"", 				"address to start http server")
	flag.BoolVar(	&cli.CacheOpt, 		"cache_opt", 	false, 				"cache option document")
	flag.BoolVar(	&cli.ServeFiles, 	"serve_files", 	true, 				"serve files from Go or not")
	flag.StringVar(	&cli.Secret, 		"secret", 		"pLsCh4nG3Th1$.AlSoThisShouldbeatLeast16bytes", "secret characters used for encryption and the like")
	flag.Parse()
}