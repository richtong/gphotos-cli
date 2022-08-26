module richtong/gphotos-cli

go 1.19

require (
	github.com/manifoldco/promptui v0.9.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.12.0
	github.com/trinhdrew1418/gphotos-cli v0.0.0-20190701062323-ed81d2e10a58
	github.com/vbauerster/mpb v3.4.0+incompatible
	golang.org/x/net v0.0.0-20220822230855-b0a4917ee28c
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094
	google.golang.org/api v0.94.0
)

require (
	cloud.google.com/go/compute v1.7.0 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

//replace google.golang.org/photoslibrary => ./vendor/google.golang.org/api/photoslibrary
// instead of replaces this with vendored point to the mirrrored one
replace github.com/trinhdrew1418/gphotos-cli => github.com/richtong/gphotos-cli v1.0.0
