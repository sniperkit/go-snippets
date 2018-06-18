Parties = new Mongo.Collection("parties");

if (Meteor.isClient) {
    angular
        .module('socially', ['angular-meteor', 'ui.router']);
    
    angular
        .module('socially').config(['$urlRouterProvider', '$stateProvider', '$locationProvider',
            function($urlRouterProvider, $stateProvider, $locationProvider) {
                $locationProvider.html5Mode(true);
                
                $stateProvider
                    .state('parties', {
                        url: '/parties',
                        templateUrl: 'parties-list.ng.html',
                        controller: 'PartiesListCtrl'
                    })
                    .state('partyDetails', {
                        url: '/parties/:partyId',
                        templateUrl: 'party-details.ng.html',
                        controller: 'PartyDetailsCtrl'
                    });
                    $urlRouterProvider.otherwise('parties');
                
            } ]);

    angular
        .module('socially').controller('PartiesListCtrl', [ '$scope', '$meteor',
            function($scope, $meteor) {
        // before go to Mongo
//                $scope.parties = [
//                    {
//                        'name': 'Dubstep-Free Zone',
//                        'description': 'Can we please just for an evening not listen to dubstep.'
//                    },
//                    {
//                        'name': 'All dubstep all the time',
//                        'description': 'Get it on!'
//                    },
//                    {
//                        'name': 'Savage lounging',
//                        'description': 'Leisure suit required. And only fiercest manners.'
//                    }
//                ];
                $scope.parties = $meteor.collection(Parties);
                $scope.remove = function(party) {
                    $scope.parties.splice( $scope.parties.indexOf(party), 1 );
                };
            } ]);
    angular.module("socially").controller("PartyDetailsCtrl", ['$scope', '$stateParams', '$meteor', 
        function($scope, $stateParams, $meteor){
            $scope.partyId = $stateParams.partyId;
        }]);
}

if (Meteor.isServer) {
    Meteor.startup(function () {
        if (Parties.find().count() === 0) {

            var parties = [
                {'name': 'Dubstep-Free Zone',
                'description': 'Can we please just for an evening not listen to dubstep.'},
                {'name': 'All dubstep all the time',
                  'description': 'Get it on!'},
                {'name': 'Savage lounging',
                  'description': 'Leisure suit required. And only fiercest manners.'}
            ];

            for (var i = 0; i < parties.length; i++)
                Parties.insert(parties[i]);

        }
    });
}
