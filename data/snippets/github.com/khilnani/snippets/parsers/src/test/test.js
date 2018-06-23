var 	jison = require("jison"),
	peg = require("pegjs");
	td = require('./testData').data,	
	log = require("dysf.utils").logger
	fs = require('fs');

var jisonFile = __dirname + '/../resources/grammar.jison';
var pegjsFile = __dirname + '/../resources/grammar.pegjs';


function start() {

	log.setLogLevel(5);

	log.info("------------------------------------------------");
	log.info("Current Dir: " + __dirname);
	log.info("------------------------------------------------");
	
	loadJison();	
}

function loadJison() {
	fs.readFile(jisonFile, 'utf8', onJisonLoad)
}

function onJisonLoad(err, data) {
	if (err) {
		log.error("Could not load: " + grammarFile);
		return log.error(err);
	}

	log.trace(data);
	log.info("*** Loaded: " + jisonFile);

	parseJison( data );

}

function parseJison( grammar ) {

	var parser = new jison.Parser(grammar);

	var parserSource = parser.generate();
	var parserResult = null;

	for(var i=0; i < td.length; i++) {
		try {
			log.info("------------------------------------------------");
			log.info("Parsing: " + td[i]);
			parserResult = parser.parse( td[i] );
			log.info("Parse result: " + parserResult);
		}
		catch( e ) {
			console.log(e);
		}
		finally {
			parserResult = null;
			log.info("------------------------------------------------");
		}
	}
	
	loadPeg();
}


function loadPeg() {
	fs.readFile(pegjsFile, 'utf8', onPegLoad);
}

function onPegLoad(err,data) {
	if (err) {
		log.error("Could not load: " + grammarFile);
		return log.error(err);
	}

	log.trace(data);
	log.info("*** Loaded: " + pegjsFile);

	parseJPeg( data );
}

function parseJPeg( grammar ) {

	var parser = peg.buildParser( grammar );

	var parserResult = null;

	for(var i=0; i < td.length; i++) {
		try {
			log.info("------------------------------------------------");
			log.info("Parsing: " + td[i]);
			parserResult = parser.parse( td[i] );
			log.info("Parse result: " + parserResult);
		}
		catch( e ) {
			console.log(e);
		}
		finally {
			parserResult = null;
			log.info("------------------------------------------------");
		}
	}
}

start();

