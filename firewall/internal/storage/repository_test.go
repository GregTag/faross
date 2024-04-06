package storage_test

import (
	"encoding/json"
	"firewall/internal/entity"
	"firewall/internal/storage"
	"firewall/internal/storage/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const jsonStr = `[
	{
		"InstanceID":"67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC",
		"auditEnabled":false,
		"format":"maven2",
		"name":"maven-snapshots",
		"namespaceConfusionProtectionEnabled":false,
		"policyCompliantComponentSelectionEnabled":false,
		"quarantineEnabled":false,
		"type":"hosted"
	},
	{
		"InstanceID":"67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC",
		"auditEnabled":false,
		"format":"maven2",
		"name":"maven-central",
		"namespaceConfusionProtectionEnabled":false,
		"policyCompliantComponentSelectionEnabled":false,
		"quarantineEnabled":false,
		"type":"proxy"
	},
	{
		"InstanceID":"67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC",
		"auditEnabled":false,
		"format":"nuget",
		"name":"nuget.org-proxy",
		"namespaceConfusionProtectionEnabled":false,
		"policyCompliantComponentSelectionEnabled":false,
		"quarantineEnabled":false,
		"type":"proxy"
	},
	{
		"InstanceID":"67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC",
		"auditEnabled":false,
		"format":"maven2",
		"name":"maven-releases",
		"namespaceConfusionProtectionEnabled":false,
		"policyCompliantComponentSelectionEnabled":false,
		"quarantineEnabled":false,
		"type":"hosted"
	},
	{
		"InstanceID":"67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC",
		"auditEnabled":true,
		"format":"pypi",
		"name":"pypi",
		"namespaceConfusionProtectionEnabled":false,
		"policyCompliantComponentSelectionEnabled":false,
		"quarantineEnabled":true,
		"type":"proxy"
	},
	{
		"InstanceID":"67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC",
		"auditEnabled":false,
		"format":"nuget",
		"name":"nuget-hosted",
		"namespaceConfusionProtectionEnabled":false,
		"policyCompliantComponentSelectionEnabled":false,
		"quarantineEnabled":false,
		"type":"hosted"
	},
	{
		"InstanceID":"67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC",
		"auditEnabled":true,
		"format":"npm",
		"name":"npm",
		"namespaceConfusionProtectionEnabled":false,
		"policyCompliantComponentSelectionEnabled":false,
		"quarantineEnabled":true,
		"type":"proxy"
	}
]`

const instance = "67B75CF1-1637A825-493619F3-2D8DC8C8-AD3DE4CC"

func TestRepositoryStore(t *testing.T) {
	db, err := driver.NewSQLiteDB("file::memory:")
	assert.NoError(t, err)

	store := storage.NewRepositoryStore(db)

	var repositories []entity.Repository
	err = json.Unmarshal([]byte(jsonStr), &repositories)
	assert.NoError(t, err)

	err = store.Save(&repositories)
	assert.NoError(t, err)

	readed, err := store.GetByInstanceId(instance)
	assert.NoError(t, err)
	assert.Equal(t, len(repositories), len(*readed))

	rep1 := &(*readed)[0]
	rep1.QuarantineEnabled = !rep1.QuarantineEnabled
	err = store.Save(rep1)
	assert.NoError(t, err)

	rep2, err := store.GetById(rep1.ID)
	assert.NoError(t, err)
	assert.Equal(t, rep1.QuarantineEnabled, rep2.QuarantineEnabled)

	future := time.Now().Add(time.Minute)
	reps, err := store.GetByInstanceIdAndSince(instance, future)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*reps))

	err = store.Save(&[]entity.Repository{})
	assert.Error(t, err)

	arr, err := store.GetByInstanceId("NO-INSTANCE-ID")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*arr))

	_, err = store.GetById(100)
	assert.Error(t, err)
}
