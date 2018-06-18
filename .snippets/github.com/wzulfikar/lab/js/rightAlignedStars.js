function rightAlignedStars (n, reversed) {
    var stars = '';
    for(var i = 0; i <= n; i++) {
        var padding = n - i;
        for (var pad = 0; pad <= padding; pad++) {
            stars += ' ';
        }
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