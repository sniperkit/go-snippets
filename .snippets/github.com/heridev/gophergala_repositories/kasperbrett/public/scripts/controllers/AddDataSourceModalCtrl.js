angular.module('kasperbrettApp').controller('AddDataSourceModalCtrl', function($scope, $http, $modalInstance) {
	$scope.dataSourceTypes = [
    	{dataSourceType: 'DsUrlScraper'}
  	];

  	$scope.selectedDataSourceType = $scope.dataSourceTypes[0];

  	$scope.preset = function(str) {
  		$scope.dataSource = {
  			name: '#Search Results on Stack Overflow for "' + str + '"',
  			interval: 30000,
  			timeout: 5000,
  			typeSettings: {
  				url: 'http://stackoverflow.com/search?q=' + str,
  				cssPath: '#mainbar > div.subheader.results-header > h2',
  				transformationScript: 'value.replace(/\\D/g, \'\')'
  			}
  		};
  	};

  	$scope.ok = function(dataSource) {
  		dataSource.type = 'DsUrlScraper';
  		dataSource.interval = parseInt(dataSource.interval);
  		dataSource.timeout = parseInt(dataSource.timeout);
  		
  		$http.post('/api/datasources', dataSource).then(function(res) {
  			var date = new Date(res.data.timestamp);
  			dataSource.id = res.data.dataSourceId;
  			dataSource.chartData = {
  				labels: [formatDateComponent(date.getHours()) + ':' + formatDateComponent(date.getMinutes()) + ':' + formatDateComponent(date.getSeconds())],
  				series: [[res.data.value]]
  			};
  			$modalInstance.close(dataSource);
  		}, function(err) {
  			console.log('test data source error', err);
  			$scope.testResponse = {
  				statusText: err.statusText,
  				status: err.status,
  				value: '',
  				error: err.data.error
  			};
  		});
  	};

  	$scope.test = function(dataSource) {
  		dataSource.type = 'DsUrlScraper';
  		dataSource.interval = parseInt(dataSource.interval);
  		dataSource.timeout = parseInt(dataSource.timeout);

  		$http.post('/api/datasources?test-only=1', dataSource).then(function(res) {
  			$scope.testResponse = {
  				statusText: res.statusText,
  				status: res.status,
  				value: res.data.value,
  				error: ''
  			};
  		}, function(err) {
  			console.log('test data source error', err);
  			$scope.testResponse = {
  				statusText: err.statusText,
  				status: err.status,
  				value: '',
  				error: err.data.error
  			};
  		});
  	};

  	function formatDateComponent(dateComponent) {
            var returnVal = dateComponent;

            if (returnVal < 10) {
                returnVal = '0' + returnVal;
            }

            return returnVal;
        }

});