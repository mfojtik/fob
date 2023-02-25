module github.com/mfojtik/fob

go 1.13

require (
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/openshift/api v0.0.0-00010101000000-000000000000
	k8s.io/api v0.0.0-20190925180651-d58b53da08f5
	k8s.io/apimachinery v0.15.7
)

replace github.com/openshift/api => github.com/openshift/api v0.0.0-20190925205819-e39b0dc4e188
