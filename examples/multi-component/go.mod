module github.com/acme.corp.sdkdemo

go 1.26.2

require github.com/vivarcus/vivarcus-sdk v0.0.0

// Local monorepo dev only. Strip replace before VPK gosdk/ deploy (ADR-21).

replace github.com/vivarcus/vivarcus-sdk => ../..
