let name: any = "foo";
let nameLen: number = (<string>name).length;
let len2: number = (name as string).length;
