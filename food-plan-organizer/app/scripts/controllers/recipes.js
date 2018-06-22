'use strict';

/**
 * @ngdoc function
 * @name foodPlanOrganizerApp.controller:MainCtrl
 * @description
 * # RecipesCtrl
 * Controller of the foodPlanOrganizerApp
 */

var ingredientHandling = function ingredientHandling($scope, Food) {
  $scope.selectedIngredient = null;
  $scope.recipe.ingredients = [];
  $scope.$watch('selectedIngredient', function(n, o) {
    if (n !== null) {
      $scope.currentIngredient = n.originalObject;
      $scope.$broadcast('add.ingredient', $scope.currentIngredient);
    }
  });
  $scope.$on('add.ingredient', function(evt, ingredient) {
    if ($scope.recipe.ingredients === undefined || $scope.recipe.ingredients === null) {
      $scope.recipe.ingredients = [];
    }
    Food.get({
      id: ingredient.id
    }, function(foodData) {
      $scope.recipe.ingredients.push({
        name: ingredient.name,
        food_id: ingredient.id,
        nutrients: foodData.nutrients
      });
    });
  });
  $scope.deleteIngredient = function(ingredient) {
    var index = $scope.recipe.ingredients.indexOf(ingredient);
    if (index !== -1) {
      $scope.recipe.ingredients.splice(index, 1);
    }
  };
};

var totalNutrients = function(recipe, nutrients, nutrientId) {
  var nutrient = null;
  for (var k = 0; k < nutrients.length; k++) {
    if (nutrients[k].id === nutrientId) {
      nutrient = nutrients[k];
    }
  }

  var total = 0.0;
  for (var i = 0; i < recipe.ingredients.length; i++) {
    var ingredient = recipe.ingredients[i];
    var multipler = 1;

    // all values are averages per 100g of editable portion
    if (ingredient.unit === 'g') {
      multipler = ingredient.volume / 100.0;
    } else if (ingredient.unit === 'cup') {
      multipler = ingredient.volume * 292.0;
    } else if (ingredient.unit === 'tbsp') {
      multipler = ingredient.volume * 18.0;
    } else if (ingredient.unit === 'tsp') {
      multipler = ingredient.volume * 6.0;
    } else {
      multipler = 1;
    }

    for (var j = 0; j < ingredient.nutrients.length; j++) {
      var nutrient = ingredient.nutrients[j];
      if (nutrient.id === nutrientId) {
        total = total + (nutrient.nutrient_value * multipler);
      }
    }
  }
  return total;
};

angular.module('foodPlanOrganizerApp')
.directive('float', function() {
  return {
    require: 'ngModel',
    link: function(scope, ele, attr, ctrl) {
      ctrl.$parsers.unshift(function(viewValue) {
        return parseFloat(viewValue);
      });
    }
  };
})
.controller('RecipesCtrl', function($scope, Recipe) {
  $scope.recipes = Recipe.query();

  $scope.truncate = function(string, n) {
    return string.length > n ? string.substr(0, n - 1) + ' …' : string;
  };
  $scope.deleteRecipe = function(recipeId) {
    Recipe.delete({
      id: recipeId
    }, function() {
      $scope.recipes = Recipe.query();
    });
  };
})
.controller('EditRecipeCtrl', function($scope, $routeParams, $location, Recipe, Nutrient, Food) {
  $scope.recipe = Recipe.get({
    id: $routeParams.id
  });
  $scope.nutrients = Nutrient.query();
  $scope.nutrientsById = function(ids) {
    var nutrients = [];
    for (var j = 0; j < ids.length; j++) {
      for (var i = 0; i < $scope.nutrients.length; i++) {
        if ($scope.nutrients[i].id === ids[j]) {
          nutrients.push($scope.nutrients[i]);
        }
      }
    }
    return nutrients;
  };
  $scope.totalNutrients = totalNutrients;
  ingredientHandling($scope, Food);

  $scope.submit = function() {
    Recipe.update($scope.recipe, function() {
      $location.path('/recipes');
    });
  };
})
.controller('NewRecipeCtrl', function($scope, $location, Recipe, Nutrient, Food) {
  $scope.recipe = {};
  $scope.nutrients = Nutrient.query();
  $scope.nutrientsById = function(ids) {
    var nutrients = [];
    for (var j = 0; j < ids.length; j++) {
      for (var i = 0; i < $scope.nutrients.length; i++) {
        if ($scope.nutrients[i].id === ids[j]) {
          nutrients.push($scope.nutrients[i]);
        }
      }
    }
    return nutrients;
  };
  $scope.totalNutrients = totalNutrients;
  ingredientHandling($scope, Food);

  $scope.submit = function() {
    Recipe.save($scope.recipe, function() {
      $location.path('/recipes');
    });
  };
});