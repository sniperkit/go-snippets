(a, b) => { let x = a + b; let y = a - b; return x + y; };
(a, b) => a + b;
(a) => { return a; };
a => { return a; };
a => a;
() => { return 0; }
() => ({a: 0});
(a, b, ...rest) => { return rest; };
(a, b = 2) => a + b;
([a, b] = [1, 2], {x: c} = {x: a + b}) => a + b + c;
