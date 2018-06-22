import { fooFunc } from "./module_declare.ts";

let a = fooFunc;

import { fooFunc as foo } from "./module_declare.ts";

let b = foo;

import * as allsyms from "./module_declare.ts";

let c = allsyms.fooFunc;
