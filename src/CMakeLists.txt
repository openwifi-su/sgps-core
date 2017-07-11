GO_GET(go_pq github.com/lib/pq)
ADD_GO_INSTALLABLE_PROGRAM(sgps-core # executable name
                           sgps-core.go # `package main` source file
                           go_pq) # everything else is a dependency
