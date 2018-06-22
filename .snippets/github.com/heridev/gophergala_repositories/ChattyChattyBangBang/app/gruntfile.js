module.exports = function(grunt) {

	// Project configuration.
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		bowercopy: {
			options: {
				srcPrefix: 'bower_components'
			},
			scripts: {
				options: {
					destPrefix: 'assets/js/vendor'
				},
				files: {
					// Target-specific file lists and/or options go here
					'framework7/framework7.js': 'framework7/dist/js/framework7.js',
					'framework7/framework7.js.map': 'framework7/dist/js/framework7.js.map',
					'framework7/framework7.min.js': 'framework7/dist/js/framework7.min.js',
					'framework7/framework7.min.js.map': 'framework7/dist/js/framework7.min.js.map'
				}
			},
			styles: {
				options: {
					destPrefix: 'assets/css'
				},
				files: {
					'framework7/framework7.min.css': 'framework7/dist/css/framework7.min.css'
				}
			}
		}
	});

	// Load the Grunt plugins.
	grunt.loadNpmTasks('grunt-bowercopy');

	// Register the default tasks.
	grunt.registerTask('move', ['bowercopy'] );
};
