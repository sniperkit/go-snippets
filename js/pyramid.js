// Print stars to create pyramid. Reversed version is included.
function pyramid (n, reversed) {
    var stars = '';
    for(var i = 0; i <= n; i++) {
        var padding = n - i;
        for (var newStar = 0; newStar < i; newStar++) {
            if (newStar <= padding) {
                stars += ' ';
            } else {
                stars += "*";
            }
        }
        stars += "\n";
    }
    if (reversed) {
        stars = stars.split("\n").reverse().join("\n");
    }
    return stars;
}