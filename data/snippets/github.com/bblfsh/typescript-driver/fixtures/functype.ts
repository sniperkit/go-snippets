type CallBack = (string) => string;

function foo(input: string): string { 
    return input; 
}
function useFoo(cb: CallBack): string { 
    return cb("bar"); 
}

let myFunc: (input: string) => string;
myFunc = foo;
