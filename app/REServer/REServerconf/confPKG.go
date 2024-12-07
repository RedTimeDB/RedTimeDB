package confer

import (
	"errors"
)

func GetNewConfer(conFileUri string) (confer Confer, err error) {

	//Parse config from yaml file
	confer.Opts, err = ParseYamlFromFile(conFileUri)

	//Check for Debug parameters: if Debug monitoring is enabled
	if confer.Opts.DebugConf.Enable {
		//If the listening address is an empty string when Enable is turned on, it indicates a configuration error
		if len(confer.Opts.DebugConf.PprofUri) == 0 {
			err = errors.New("When the Enable parameter in Debug is True, the pprof_uri parameter cannot be empty, and a meaningful pprof listening address must be filled in.")
			return
		}
	}

	//If caching is enabled
	if confer.Opts.CacheConf.Enable {
		//Cache related configuration parameters security check: set the minimum number of items in the cache
		if confer.Opts.CacheConf.MaxItemsSize < 1024 {
			confer.Opts.CacheConf.MaxItemsSize = 1024
		}
	}

	//cache default expiration ms lower limit check
	if confer.Opts.CacheConf.DefaultExpiration <= 0 {
		confer.Opts.CacheConf.DefaultExpiration = 0
	}

	//Cache expired kv automatic cleaning cycle Minimum cycle limit (minimum 10s）
	if confer.Opts.CacheConf.CleanupInterval <= 10 {
		confer.Opts.CacheConf.CleanupInterval = 10
	}

	return
}
