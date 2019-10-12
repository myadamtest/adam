package mygen

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func SqlExec(db *sql.DB, sqlstr string, args ...interface{}) (int64, int64, error) {
	result, err := db.Exec(sqlstr, args...)
	if err != nil {
		return 0, 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	num, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return id, num, nil
}

func SqlSelect(db *sql.DB, sqlstr string, args ...interface{}) ([]map[string]string, error) {
	rows, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	var rets = make([]map[string]string, 0)

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var ret = make(map[string]string)
		var value string
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			ret[columns[i]] = value
		}

		rets = append(rets, ret)
	}

	return rets, err
}

func HttpGet(cli *http.Client, reqUrl string, v interface{}) error {
	resp, err := cli.Get(reqUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

func HttpPostForm(cli *http.Client, reqUrl string, values url.Values, v interface{}) error {
	resp, err := cli.PostForm(reqUrl, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

func HttpPostJson(cli *http.Client, reqUrl string, values []byte, v interface{}) error {
	resp, err := cli.Post(reqUrl, "application/json", bytes.NewReader(values))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

func HttpPostJsonNew(cli *http.Client, reqUrl string, values []byte, v interface{}) error {
	resp, err := cli.Post(reqUrl, "application/json", bytes.NewReader(values))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("http post json fail, status:%d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}
