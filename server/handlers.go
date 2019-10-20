package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type DataRow struct {

}
type DataRowSlice []DataRow

func (serv *Server) GetData(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(path.Join(serv.conf.Assets, "data.json"))
	if err != nil {
		if err == os.ErrNotExist {
			serv.lg.Infof("data file doesn't exist: %s", err)
			serv.JSONResp(w, DataRowSlice{})
			return
		}
		serv.lg.Errorf("can't read file: %s", err)
		serv.JSONResp(w, DataRowSlice{})
		return
	}
	if _, err := w.Write(data); err != nil {
		serv.lg.Errorf("can't write data: %s", err)
		serv.JSONResp(w, DataRowSlice{})
		return
	}
}

func (serv *Server) PostData(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serv.lg.Infof("data file doesn't exist: %s", err)
		return
	}
	if err := ioutil.WriteFile(path.Join(serv.conf.Assets, "data.json"), data, os.ModePerm); err != nil {
		serv.lg.Errorf("can't save data: %s", err)
		return
	}
}
