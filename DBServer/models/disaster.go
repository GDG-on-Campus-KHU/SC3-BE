package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Disaster struct {
	SN        int    `json:"SN"`
	Timestamp string `json:"CRT_DT"`
	Message   string `json:"MSG_CN"`
	Region    string `json:"RCPTN_RGN_NM"`
	Step      string `json:"EMRG_STEP_NM"`
	Segment   string `json:"DST_SE_NM"`
	RegDate   string `json:"REG_YMD"`
	ModDate   string `json:"MDFCN_YMD"`
}

type Person struct {
	SN        string             `bson:"sn"`
	Timestamp primitive.DateTime `bson:"timestamp"`
	Region    []string           `bson:"region"`
	Missing   Missing            `bson:"missing"`
	Accessed  bool               `bson:"accessed"`
}

type Missing struct {
	Name        string `bson:"name"`
	Age         int    `bson:"age"`
	Sex         string `bson:"sex"`
	Description string `bson:"description"`
}
