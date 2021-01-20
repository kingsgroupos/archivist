/*
Owned by watcher. DON'T change me.
*/

var argv = process.argv.slice(2)
if (argv.length != 1) {
	console.error('Usage: node js2json.js file.js')
	process.exit(1)
}

try {
	var data = require(argv[0])
	var json = JSON.stringify(data, null, 4)
	console.log(json)
} catch (err) {
	console.error('Failed to load "' + argv[0] + '"')
	console.error(err.stack)
	process.exit(1)
}
