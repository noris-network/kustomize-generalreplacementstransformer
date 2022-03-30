To see the transformed output, `cd` to this directory and run

    kustomize build --enable-alpha-plugins

Sometimes it might be useful to see debug `GeneralReplacementsTransformer`'s
debug output. When `GeneralReplacementsTransformer` is in your `PATH`, this can be
achieved, e.g., by running

    GeneralReplacementsTransformer transformer.yaml < <(for f in *.yml; do cat $f; echo "---"; done) | grep ^##

This will print out

    ## GeneralReplacementsTransformer: namespace = "myapp"
    ## GeneralReplacementsTransformer: username = "testuser"
    ## GeneralReplacementsTransformer: db.name = "mydb"
    ## GeneralReplacementsTransformer: db.host = "mongodb"
    ## GeneralReplacementsTransformer: db.port = "27017"
    ## GeneralReplacementsTransformer: cacheHost = "redis"
    ## GeneralReplacementsTransformer: cachePort = "6379"
    ## GeneralReplacementsTransformer: cacheDb = "3"
    ## GeneralReplacementsTransformer: password = "GeneralReplacementsTransformer"
    ## GeneralReplacementsTransformer: select value "oops" not found
