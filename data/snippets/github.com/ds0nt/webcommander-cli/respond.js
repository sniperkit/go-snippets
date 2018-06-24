#!/bin/node

function fmpInsight (wordCount) {
    var output = "", index = old = Math.floor(Math.random()*30), word="", related = null,
        associations = [{data:"that's basically it.", related:[-1]},
                        {data:"i want to start to put",related:[12, 19, 20, 22, 23]},
                        {data:"we need to start putting",related:[12, 19, 20, 22, 23]},
                        {data:"we need to start putting more",related:[12, 19, 20, 22, 23]},
                        {data:"do like ",related:[16, 19, 30, -1]},
                        {data:"I've sent you an email.",related:[-1]},
                        {data:"we need to have more money coming in",related:[4 -1]},
                        {data:"My trainer says",related:[6, 9, 14, 16, 17 ]},
                        {data:"first, we need to",related:[4, 28]},
                        {data:"I need some coffee.",related:[-1]},
                        {data:"Where's my phone!?",related:[20, -1]},
                        {data:"My phone! Now!",related:[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
                        {data:"cookies",related:[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
                        {data:"In Portugal",related:[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
                        {data:"there is a saying in Portugal",related:[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
                        {data:"So that's the take",related:[-1]},
                        {data:"Red Train",related:[18]},
                        {data:"iService",related:[18]},
                        {data:"Those bastards!",related:[-1]},
                        {data:"actionable",related:[0]},
                        {data:"fast",related:[21, -1]},
                        {data:"reliable",related:[22, -1]},
                        {data:"affordable",related:[23, -1]},
                        {data:"Maria.",related:[24, -1]},
                        {data:"Maria?",related:[25, -1]},
                        {data:"Maria!",related:[11, -1]},
                        {data:"Let me ask you a question",related:[-1, ]},
                        {data:"design",related:[28]},
                        {data:"build",related:[29]},
                        {data:"promote",related:[-1]},
                        {data:"Vital",related:[0]}];
    while (wordCount-- >0) {
        if (index == -1) {
            old = index = Math.floor(Math.random()*30);
            output += ". ";
        }
        output += associations[index].data + " ";
        related = associations[old].related;
        index = related[Math.floor(Math.random()*related.length)];
        old = index;
    }
    return output;
}

console.log(fmpInsight(50))
