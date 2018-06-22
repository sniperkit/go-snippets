module.exports = function(grunt) {
	'use strict';

	// includes
	grunt.loadNpmTasks('grunt-bower-task');
	grunt.loadNpmTasks('grunt-contrib-clean');
	grunt.loadNpmTasks('grunt-contrib-copy');
	grunt.loadNpmTasks('grunt-contrib-csslint');
	grunt.loadNpmTasks('grunt-contrib-cssmin');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	// Project configuration.
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
			clean: {
				bower: ["public/vendor/"],
				cssvendor: ["public/css/vendor/"],
				cssdist: ["public/css/dist/"],
				jsvendor: ["public/js/vendor/"],
			},
			// ========== CSS ==========
			csslint: {
					src: ['css/km.css'],
					options: {
						ids: false
					}
				},
				cssmin: {
					dist: {
						options: {
							report: "gzip",
						},
						files: {
							"public/css/dist/app.min.css": ['public/css/vendor/*css']
						}
					}
				},
        // ========== END CSS ==========

        // ========== JS ==========

        jshint: {
					files: ["js/km.js"],
        	options: {
                // Development.
                "debug"         : false,  // Allow debugger statements e.g. browser breakpoints.
                "devel"         : false,   // Allow developments statements e.g. `console.log();`.

                // Settings
                "passfail"      : false,  // Stop on first error.
                "maxerr"        : 100,    // Maximum error before stopping.

                // Predefined globals whom JSHint will ignore.
                "browser"       : true,   // Standard browser globals e.g. `window`, `document`.
                "node"          : true,
                "rhino"         : false,
                "couch"         : false,
                "wsh"           : false,   // Windows Scripting Host.
                "jquery"        : true,
                "prototypejs"   : false,
                "mootools"      : false,
                "dojo"          : false,

                // The Good Parts.
                "asi"           : false,  // Tolerate Automatic Semicolon Insertion (no semicolons).
                "laxbreak"      : true,   // Tolerate unsafe line breaks e.g. `return [\n] x` without semicolons.
                "bitwise"       : true,   // Prohibit bitwise operators (&, |, ^, etc.).
                "boss"          : false,  // Tolerate assignments inside if, for & while. Usually conditions & loops are for comparison, not assignments.
                "curly"         : false,   // Require {} for every new block or scope.
                "eqeqeq"        : true,   // Require triple equals i.e. `===`.
                "eqnull"        : false,  // Tolerate use of `== null`.
                "evil"          : false,  // Tolerate use of `eval`.
                "expr"          : false,  // Tolerate `ExpressionStatement` as Programs.
                "forin"         : false,  // Tolerate `for in` loops without `hasOwnPrototype`.
                "immed"         : true,   // Require immediate invocations to be wrapped in parens e.g. `( function(){}() );`
                "latedef"       : true,   // Prohibit variable use before definition.
                "loopfunc"      : false,  // Allow functions to be defined within loops.
                "noarg"         : true,   // Prohibit use of `arguments.caller` and `arguments.callee`.
                "regexp"        : true,   // Prohibit `.` and `[^...]` in regular expressions.
                "regexdash"     : false,  // Tolerate unescaped last dash i.e. `[-...]`.
                "scripturl"     : true,   // Tolerate script-targeted URLs.
                "shadow"        : false,  // Allows re-define variables later in code e.g. `var x=1; x=2;`.
                "supernew"      : false,  // Tolerate `new function () { ... };` and `new Object;`.
                "undef"         : true,   // Require all non-global variables be declared before they are used.
                "es5"           : false,   // If ES5 syntax should be allowed.
                "strict"        : false,  // Require the "use strict"; pragma.
                "onecase"       : true,

                // Personal styling preferences.
                "newcap"        : true,   // Require capitalization of all constructor functions e.g. `new F()`.
                "noempty"       : true,   // Prohibit use of empty blocks.
                "nonew"         : true,   // Prohibit use of constructors for side-effects.
                "nomen"         : true,   // Prohibit use of initial or trailing underbars in names.
                "onevar"        : false,  // Allow only one `var` statement per function.
                "plusplus"      : false,  // Prohibit use of `++` & `--`.
                "sub"           : true,  // Tolerate all forms of subscript notation besides dot notation e.g. `dict['key']` instead of `dict.key`.
                "trailing"      : false,   // Prohibit trailing whitespaces.
                "white"         : false,   // Check against strict whitespace and indentation rules.

                // globals
                globals: {
                  "curr_img": true,
                  "uri": true,
									"images": true
                }
            }
        },

				bower: {
					install: {
						options: {
							cleanBowerDir: true,
							layout: 'byType',
							targetDir: 'public/vendor/',
							verbose: true
						}
					}
				},
				copy: {
					jsvendor: {
						flatten: true,
						expand: true,
						src: ['public/vendor/requirejs/require.js', 'public/vendor/jquery/jquery.js', 'public/vendor/bootstrap/bootstrap.js', 'public/vendor/swig/js/swig.js'],
						dest: 'public/js/vendor/'
					},
					cssvendor: {
						flatten: true,
						expand: true,
						src: ['public/vendor/bootstrap/bootstrap.css'],
						dest: 'public/css/vendor/'
					},
				}

				// ========== END JS ==========
    });

    // Default task: the works
		grunt.registerTask('default', ['js', 'css']);
		// Make it
		grunt.registerTask('make', ['clean', 'cssmin']);
		// test only
		grunt.registerTask('test', ['csslint', 'jshint']);
		// JS
		grunt.registerTask('js', ["clean:jsvendor", "bower", "copy:jsvendor", "clean:bower"]);
		// CSS
		grunt.registerTask('css', ["clean:cssvendor", "clean:cssdist", "bower", "copy:cssvendor", "csslint", "cssmin"]);
};
