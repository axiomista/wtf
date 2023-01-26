package oura

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testSettings = Settings{
		accessToken: "MFKOR3Z4WHJ7XBHURYSJNSZ7DMCYH7P2",
		myName:      "Me",
		days:        4,
	}
)

func Test_getUserInfo(t *testing.T) {
	testClient := NewClient(&testSettings)
	infoResult, err := testClient.getUserInfo()
	assert.Equal(t, nil, err)
	assert.Equal(t, "vanessanicole@gmail.com", infoResult.Email)
}

func Test_getSleeps(t *testing.T) {
	testClient := NewClient(&testSettings)
	testClient.start = "2021-09-05"
	testClient.end = "2021-09-07"
	sleepResult, err := testClient.getSleeps()
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, sleepResult)
}
