package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	INN         string `json:"INN"`
	KPP         string `json:"KPP"`
	Leader      string `json:"Leader"`
	CompanyName string `json:"Company_name"`
}

func TestGoParse(t *testing.T) {
	var testData testData
	testData.INN = "1657201932"
	testData.KPP = "165701001"
	testData.CompanyName = "ОБЩЕСТВО С ОГРАНИЧЕННОЙ ОТВЕТСТВЕННОСТЬЮ \"ОНЛАЙН КОНСАЛТ\""
	testData.Leader = "Матвеев Михаил Анатольевич"

	data := GoParser("1657201932")
	assert.Equal(t, data.INN, testData.INN)
	assert.Equal(t, data.KPP, testData.KPP)
	assert.Equal(t, data.CompanyName, testData.CompanyName)
	assert.Equal(t, data.Leader, testData.Leader)

}
