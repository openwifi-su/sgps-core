GO_GET(go_pq github.com/lib/pq)
GO_GET(go_viper github.com/spf13/viper)
ADD_GO_INSTALLABLE_PROGRAM(sgps-core # executable name
                           sgps-core.go # `package main` source file
                           go_pq
                           go_viper) # everything else is a dependency
