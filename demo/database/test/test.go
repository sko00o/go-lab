package test

var DSNTable = map[string]string{
	//"mymysql": "tcp:localhost:13306*demo_test/root/toor",
	//"mysql":    "root:toor@tcp(localhost:13306)/demo_test?parseTime=True",
	"postgres": "dbname=demo_test user=root password=toor port=15432 sslmode=disable",
}
