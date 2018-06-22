var str = "good!";
var str1 = "str1!good!";
var str2 = "good!str2!";
var dmp = new diff_match_patch();

var patchs = dmp.patch_make(str, str1);
// str
var outputs = dmp.patch_apply(patchs, str)
str = outputs[0];
result = outputs[1];
for (var i=0;i<result.length;i++){
    if (!result[i]){
        console.log("result:"+result);
    }
}
// str2
outputs = dmp.patch_apply(patchs, str2)
str2 = outputs[0];
result = outputs[1];
for (var i=0;i<result.length;i++){
    if (!result[i]){
        console.log("result:"+result);
    }
}

patchs = dmp.patch_make(str, str2);
// str
outputs = dmp.patch_apply(patchs, str)
str = outputs[0];
result = outputs[1];
for (var i=0;i<result.length;i++){
    if (!result[i]){
        console.log("result:"+result);
    }
}
// str1
outputs = dmp.patch_apply(patchs, str1)
str1 = outputs[0];
result = outputs[1];
for (var i=0;i<result.length;i++){
    if (!result[i]){
        console.log("result:"+result);
    }
}

console.log(str);