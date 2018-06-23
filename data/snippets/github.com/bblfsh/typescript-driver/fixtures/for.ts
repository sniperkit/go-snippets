let arr = [1, 2, 3];

for (let i=0; i < arr.length; i++) {
    console.log(i);
    break;
}

for (let item of arr) {
    console.log(item);
    continue;
}

for (let index in arr) {
    console.log(index);
}
