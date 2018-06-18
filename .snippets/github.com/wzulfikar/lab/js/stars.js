// print stars :v
function stars (n, reversed) {
    var stars = '';
    for(var i = 0; i <= n; i++) {
        for (var newStar = 0; newStar < i; newStar++) {
            stars += "*";            
        }
        stars += "\n";
    }
    if (reversed) {
        stars = stars.split("\n").reverse().join("\n");
    }
    return stars;
}

const N = 2;
console.log("Number of stars:", N);

console.log(star(N));

console.log('Reversed version:')
console.log(star(N, true));