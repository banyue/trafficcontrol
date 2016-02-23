// Copyright 2015 Comcast Cable Communications Management, LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file was initially generated by gen_to_start.go (add link), as a start
// of the Traffic Ops golang data model

package api

import (
	"encoding/json"
	_ "github.com/Comcast/traffic_control/traffic_ops/experimental/server/output_format" // needed for swagger
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type StatsSummary struct {
	Id                  int64             `db:"id" json:"id"`
	CdnName             string            `db:"cdn_name" json:"cdnName"`
	DeliveryserviceName string            `db:"deliveryservice_name" json:"deliveryserviceName"`
	StatName            string            `db:"stat_name" json:"statName"`
	StatValue           float64           `db:"stat_value" json:"statValue"`
	SummaryTime         time.Time         `db:"summary_time" json:"summaryTime"`
	StatDate            time.Time         `db:"stat_date" json:"statDate"`
	Links               StatsSummaryLinks `json:"_links" db:-`
}

type StatsSummaryLinks struct {
	Self string `db:"self" json:"_self"`
}

// @Title getStatsSummaryById
// @Description retrieves the stats_summary information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    StatsSummary
// @Resource /api/2.0
// @Router /api/2.0/stats_summary/{id} [get]
func getStatsSummaryById(id int, db *sqlx.DB) (interface{}, error) {
	ret := []StatsSummary{}
	arg := StatsSummary{}
	arg.Id = int64(id)
	queryStr := "select *, concat('" + API_PATH + "stats_summary/', id) as self "
	queryStr += " from stats_summary where id=:id"
	nstmt, err := db.PrepareNamed(queryStr)
	err = nstmt.Select(&ret, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nstmt.Close()
	return ret, nil
}

// @Title getStatsSummarys
// @Description retrieves the stats_summary
// @Accept  application/json
// @Success 200 {array}    StatsSummary
// @Resource /api/2.0
// @Router /api/2.0/stats_summary [get]
func getStatsSummarys(db *sqlx.DB) (interface{}, error) {
	ret := []StatsSummary{}
	queryStr := "select *, concat('" + API_PATH + "stats_summary/', id) as self "
	queryStr += " from stats_summary"
	err := db.Select(&ret, queryStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ret, nil
}

// @Title postStatsSummary
// @Description enter a new stats_summary
// @Accept  application/json
// @Param                 Body body     StatsSummary   true "StatsSummary object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/stats_summary [post]
func postStatsSummary(payload []byte, db *sqlx.DB) (interface{}, error) {
	var v StatsSummary
	err := json.Unmarshal(payload, &v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "INSERT INTO stats_summary("
	sqlString += "cdn_name"
	sqlString += ",deliveryservice_name"
	sqlString += ",stat_name"
	sqlString += ",stat_value"
	sqlString += ",summary_time"
	sqlString += ",stat_date"
	sqlString += ") VALUES ("
	sqlString += ":cdn_name"
	sqlString += ",:deliveryservice_name"
	sqlString += ",:stat_name"
	sqlString += ",:stat_value"
	sqlString += ",:summary_time"
	sqlString += ",:stat_date"
	sqlString += ")"
	result, err := db.NamedExec(sqlString, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title putStatsSummary
// @Description modify an existing stats_summaryentry
// @Accept  application/json
// @Param   id              path    int     true        "The row id"
// @Param                 Body body     StatsSummary   true "StatsSummary object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/stats_summary/{id}  [put]
func putStatsSummary(id int, payload []byte, db *sqlx.DB) (interface{}, error) {
	var v StatsSummary
	err := json.Unmarshal(payload, &v)
	v.Id = int64(id) // overwrite the id in the payload
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "UPDATE stats_summary SET "
	sqlString += "cdn_name = :cdn_name"
	sqlString += ",deliveryservice_name = :deliveryservice_name"
	sqlString += ",stat_name = :stat_name"
	sqlString += ",stat_value = :stat_value"
	sqlString += ",summary_time = :summary_time"
	sqlString += ",stat_date = :stat_date"
	sqlString += " WHERE id=:id"
	result, err := db.NamedExec(sqlString, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title delStatsSummaryById
// @Description deletes stats_summary information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    StatsSummary
// @Resource /api/2.0
// @Router /api/2.0/stats_summary/{id} [delete]
func delStatsSummary(id int, db *sqlx.DB) (interface{}, error) {
	arg := StatsSummary{}
	arg.Id = int64(id)
	result, err := db.NamedExec("DELETE FROM stats_summary WHERE id=:id", arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}
