var name = 'zach';

while (true) {
    var name = 'obma';
    alert(name);
    break;
}

alert(name);

let name1 = 'zach'

while (true) {
    let name1 = 'obama';
    alert(name1);
    break;
}

alert(name1);


const PI = Math.PI;

alert(PI);

class Animal {
    constructor(){
        this.type = 'animal'
    }
    says(say){
        alert(this.type + ' says ' + say)
    }
}

let animal = new Animal()
    animal.says('hello') //animal says hello

    class Cat extends Animal {
        constructor(){
            super()
                this.type = 'cat'
        }
    }

let cat = new Cat()
    cat.says('hello')


