To see the transformed output, `cd` to this directory and run

    kustomize build --enable-alpha-plugins

Sometimes it's useful to see `GeneralReplacementsTransformer`'s debug output.
When `GeneralReplacementsTransformer` is in your `PATH`, this can be achieved
by commenting out `- transformer.yaml` in `kustomization.yaml` and running

    kustomize build | GeneralReplacementsTransformer transformer.yaml | grep ^#

This will print out

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
