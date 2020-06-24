package badgerds

import (
	"os"

	"github.com/ipfs/go-ipfs/plugin"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/threefoldtech/0-stor/client/datastor"

	"github.com/zaibon/zstords"
)

// Plugins is exported list of plugins that will be loaded
var Plugins = []plugin.Plugin{
	&zstorPlugin{},
}

type zstorPlugin struct{}

var _ plugin.PluginDatastore = (*zstorPlugin)(nil)

func (*zstorPlugin) Name() string {
	return "ds-ztsor"
}

func (*zstorPlugin) Version() string {
	return "0.1.0"
}

func (*zstorPlugin) Init(_ *plugin.Environment) error {
	return nil
}

func (*zstorPlugin) DatastoreTypeName() string {
	return "zstords"
}

type datastoreConfig struct {
	zstords.Options
}

// BadgerdsDatastoreConfig returns a configuration stub for a badger datastore
// from the given parameters
func (*zstorPlugin) DatastoreConfigParser() fsrepo.ConfigFromMap {
	return func(params map[string]interface{}) (fsrepo.DatastoreConfig, error) {
		var c = datastoreConfig{}
		// var ok bool
		c.JobCount = 1
		c.Config.DataStor.Shards = []datastor.ShardConfig{
			{
				Address:   "localhost:9001",
				Namespace: "default",
				Password:  "",
			},
		}
		c.Config.DataStor.Pipeline.BlockSize = 1024 * 1024 // 1MiB
		c.MetaPath = "ipfs-meta"

		if err := os.MkdirAll(c.MetaPath, 0770); err != nil {
			return nil, err
		}

		return &c, nil
	}
}

func (c *datastoreConfig) DiskSpec() fsrepo.DiskSpec {
	return map[string]interface{}{
		"type": "zstords",
		"path": c.MetaPath,
	}
}

func (c *datastoreConfig) Create(path string) (repo.Datastore, error) {
	return zstords.NewDatastore(&c.Options)
}
