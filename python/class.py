class Student(object):
    def __init__(self, name, score):
        self.name = name
        self.score = score
    
    def print_score(self):
        print("%s: %s" % (self.name, self.score))


a = Student("Foo", 56)
b = Student("Bar", 110)

a.print_score()
b.print_score()

class Animal(object):
    def run(self):
        print("Animal is running....")

class Dog(Animal):
    def run(self):
        print("Dog is running...")

class Cat(Animal):
    def run(self):
        print("Cat is running...")

def run_twice(animal):
    animal.run()
    animal.run()

d = Dog()
d.run()
c = Cat()
c.run()

run_twice(Animal())
run_twice(Dog())
run_twice(Cat())

