package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/niftynei/glightning/jrpc2"
	"log"
	"reflect"
)

type ChannelActivity struct {
	Short_Channel_id string `json:"short_channel_id"`
	Actions          int    `json:"count"`
	Msatoshi         int64  `json:"msatoshi"`
	Direction        string `json:"direction"`
}

func (f *ChannelActivity) New() interface{} {
	return &ChannelActivity{}
}

func (f *ChannelActivity) Name() string {
	return "channel_activity"
}

func (z *ChannelActivity) Call() (jrpc2.Result, error) {
	return activitySummary()
}

func activitySummary() (interface{}, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Printf("db open failed: %s\n", err.Error())
		return nil, err
	}

	rows, err := db.Query(
		`SELECT 
			c.short_channel_id, 
			count(h.msatoshi) actions,
			sum(h.msatoshi), 
      CASE h.direction
        WHEN 1 THEN 'send'
        WHEN 0 THEN 'receive'
      END
		FROM channels c
		JOIN channel_htlcs h ON c.id=h.channel_id
		WHERE failuremsg IS NULL
		GROUP BY c.short_channel_id, h.direction;`)

	if err != nil {
		log.Printf("activity query failed: %s\n", err.Error())
		return nil, err
	}

	result := make([]ChannelActivity, 0)

	for rows.Next() {
		a := &ChannelActivity{}
		err := scanToStruct(a, rows)
		result = append(result, *a)
		if err != nil {
			log.Printf("db query fields error: %s", err.Error())
			return nil, err
		}
	}

	db.Close()
	return result, nil

}

func scanToStruct(obj interface{}, rows *sql.Rows) error {
	s := reflect.ValueOf(obj).Elem()
	fields := make([]interface{}, 0)
	for i := 0; i < s.NumField(); i++ {
		var f interface{}
		fields = append(fields, &f)
	}

	err := rows.Scan(fields...)

	for i := 0; i < s.NumField(); i++ {
		var raw_value = *fields[i].(*interface{})
		setFieldValue(s.Field(i), raw_value)
	}

	return err
}

func setFieldValue(field reflect.Value, val interface{}) {
	if val == nil {
		return
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(val.(string))
	case reflect.Int, reflect.Int64:
		field.SetInt(val.(int64))
	}

}
