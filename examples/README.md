To see the transformed output, `cd` to this directory and run

    kustomize build --enable-alpha-plugins

To see some log output that might be helpful for debugging set `GRT_LOG=/path/to/logfile`
in the environment, e.g.run

    GRT_LOG=/tmp/grt.log kustomize build --enable-alpha-plugins

    cat /tmp/grt.log
    ## GeneralReplacementsTransformer: namespace = "demo"
    ## GeneralReplacementsTransformer: username = "testuser"
    ## GeneralReplacementsTransformer: db.host = "mongodb"
    ## GeneralReplacementsTransformer: db.name = "mydb"
    ## GeneralReplacementsTransformer: db.port = "27017"
    ## GeneralReplacementsTransformer: dummy.cachePort = "6379"
    ## GeneralReplacementsTransformer: dummy.cacheDb = "3"
    ## GeneralReplacementsTransformer: dummy.cacheHost = "redis"
    ## GeneralReplacementsTransformer: password = "s3cr3t1234"
    ## GeneralReplacementsTransformer: select value "oops" not found
