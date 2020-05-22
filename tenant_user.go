package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func setupDBConnection() string {
	port := 5432
	if p := os.Getenv("DB_PORT"); p != "" {
		i, err := strconv.ParseInt(p, 10, 32)
		if err == nil {
			port = int(i)
		} else {
			panic("cannot convert DB_PORT to integer")
		}
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "clm-aus-t3vfey.bmc.com"
		//panic("failed to get env DB_HOST")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		//panic("failed to get env DB_USER")
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		//panic("failed to get env DB_PASSWORD")
		password = "password"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		//panic("failed to get env DB_NAME")
		dbname = "panama"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return psqlInfo

}

func generateTenantInfo() {

	if h := os.Getenv("INTERVAL_MINUTES"); h != "" {
		i, err := strconv.ParseInt(h, 10, 32)
		if err == nil {
			tenantInterval = time.Duration(i) * time.Minute
			fmt.Printf("set metrics querying cycle to %d minutes \n", i)
		} else {
			panic("cannot convert INTERVAL_MINUTES to integer")
		}
	}

	for {

		// open the db
		psgInfo := setupDBConnection()
		db, err := sql.Open("postgres", psgInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		t := time.Now()

		// time in milli seconds
		unixTime := t.Unix() * 1000
		currentTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		fmt.Printf("current time is %s, ts is %d\n", currentTime, unixTime)

		rows, err := db.Query(
			`SELECT t.name,  MAX(f.updated_at) FROM public."Flows" f 
			 INNER JOIN public."Tenants" t on f.tenant_id=t.id
			 GROUP BY t.name`)

		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var name string
			var lastSeen time.Time
			if e := rows.Scan(&name, &lastSeen); e == nil {
				fmt.Printf("row:  %s %s\n", name, lastSeen.String())
				tenantUsage.WithLabelValues("DC-1", name, currentTime).Set(float64(lastSeen.Unix()))
			} else {
				panic(e)
			}
		}
		time.Sleep(time.Duration(tenantInterval) * time.Minute)
	}
}
