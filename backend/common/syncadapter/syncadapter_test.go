package syncadapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SyncAdapterTestSuite struct {
	suite.Suite
}

func (suite *SyncAdapterTestSuite) TestNewFieldMapper() {
	mapper := NewFieldMapper()
	assert.NotNil(suite.T(), mapper)
}

func (suite *SyncAdapterTestSuite) TestNewSyncMetrics() {
	metrics := NewSyncMetrics()
	assert.NotNil(suite.T(), metrics)
}

func (suite *SyncAdapterTestSuite) TestExternalSyncConfig() {
	config := ExternalSyncConfig{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "test",
		Password: "password",
		Database: "testdb",
		Charset:  "utf8mb4",
	}

	assert.Equal(suite.T(), "mysql", config.Type)
	assert.Equal(suite.T(), "localhost", config.Host)
	assert.Equal(suite.T(), 3306, config.Port)
}

func (suite *SyncAdapterTestSuite) TestSyncTask() {
	task := SyncTask{
		TaskId:   "test-task-001",
		Type:     "order",
		OrderId:  123,
		Attempts: 0,
	}

	assert.Equal(suite.T(), "test-task-001", task.TaskId)
	assert.Equal(suite.T(), "order", task.Type)
	assert.Equal(suite.T(), int64(123), task.OrderId)
}

func TestSyncAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(SyncAdapterTestSuite))
}
