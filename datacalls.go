package GoSDK

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	_DATA_PREAMBLE = "/api/v/1/data/"
)

func (u *UserClient) InsertData(collection_id string, data interface{}) error {
	return insertdata(u, collection_id, data)
}

func (d *DevClient) InsertData(collection_id string, data interface{}) error {
	return insertdata(d, collection_id, data)
}

func insertdata(c cbClient, collection_id string, data interface{}) error {
	resp, err := post(_DATA_PREAMBLE+collection_id, data, c.creds())
	if err != nil {
		return fmt.Errorf("Error inserting: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error inserting: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) GetData(collection_id string, query [][]map[string]interface{}) (map[string]interface{}, error) {
	return getdata(u, collection_id, query)
}

func (d *DevClient) GetData(collection_id string, query [][]map[string]interface{}) (map[string]interface{}, error) {
	return getdata(d, collection_id, query)
}

func getdata(c cbClient, collection_id string, query [][]map[string]interface{}) (map[string]interface{}, error) {
	var qry map[string]string
	if query != nil {
		b, jsonErr := json.Marshal(query)
		if jsonErr != nil {
			return nil, fmt.Errorf("JSON Encoding error: %v", jsonErr)
		}
		qryStr := url.QueryEscape(string(b))
		qry = map[string]string{"query": qryStr}
	} else {
		qry = nil
	}
	resp, err := get(_DATA_PREAMBLE+collection_id, qry, c.creds())
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting data: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) UpdateData(collection_id string, query [][]map[string]interface{}, changes map[string]interface{}) error {
	return getdata(u, collection_id, query, changes)
}

func (d *DevClient) UpdateData(collection_id string, query [][]map[string]interface{}, changes map[string]interface{}) error {
	return getdata(d, collection_id, query)
}

func updatedata(c cbClient, collection_id string, query [][]map[string]interface{}, changes map[string]interface{}) error {
	body := map[string]interface{}{
		"query": query,
		"$set":  changes,
	}
	resp, err := put(_DATA_PREAMBLE+collection_id, body, c.creds())
	if err != nil {
		return fmt.Errorf("Error updating data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating data: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) DeleteData(collection_id string, query [][]map[string]interface{}) (map[string]interface{}, error) {
	return deletedata(u, collection_id, query)
}

func (d *DevClient) DeleteData(collection_id string, query [][]map[string]interface{}) (map[string]interface{}, error) {
	return deletedata(d, collection_id, query)
}

func deletedata(c cbClient, collection_id string, query [][]map[string]interface{}) error {
	var qry map[string]string
	if query != nil {
		b, jsonErr := json.Marshal(query)
		if jsonErr != nil {
			return fmt.Errorf("JSON Encoding error: %v", jsonErr)
		}
		qryStr := url.QueryEscape(string(b))
		qry = map[string]string{"query": qryStr}
	} else {
		return fmt.Errorf("Must supply a query to delete")
	}
	resp, err := delete(_DATA_PREAMBLE+collection_id, qry, c.creds())
	if err != nil {
		return fmt.Errorf("Error deleting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting data: %v", resp.Body)
	}
	return nil
}
