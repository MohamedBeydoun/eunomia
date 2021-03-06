package tpl

// ConsoleTemplate ...
func ConsoleTemplate() []byte {
	return []byte(`const repl = require("repl");
global.mongoose = require("mongoose")
global.App =  require("./dist/app.js").app

global.mongoose.connect("{{ .DBURL }}", {
    useNewUrlParser: true,
    useFindAndModify: false,
    useUnifiedTopology: true,
    useCreateIndex: true
});
global.mongoose.set("useCreateIndex", true);

{{ range $model := .Models }}
global.{{ $model | ToTitle }} = require("./dist/database/interactions/{{ $model }}.js").{{ $model }}DBInteractions
{{ end }}

r = repl.start({prompt: "atlas> ", useGlobal: true, experimentalReplAwait: true});
`)
}
