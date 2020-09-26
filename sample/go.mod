module github.com/oncilla/boa/sample

go 1.13

replace github.com/oncilla/boa => ../../boa

require (
	github.com/mitchellh/mapstructure v1.1.2
	github.com/oncilla/boa v0.1.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
