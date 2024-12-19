# A simple scripting language

Evie is a personal project of a simple interpreted language with features like structures, dictionaries, modules, and more.

## Variables

```
var x = 20
x = "Now x is a text!"
y = 3      // PANIC: "y" does not exists
var x = 2  // PANIC: can not redeclare var "x"
```

## Functions
```
fn add(n1, n2){
  return n1 + n2
}

var result = add(2,4)
print(result) // 6
```

## Structures
```
struct Person{
  name,
  age
}

var person1 = Person{}
person1.name = "John"
person1.age = 25
person1.oops = "oops" // PANIC: person struct does not have oops property

var person2 = Person{
  name: "Pedro",
  age: 24,
  realProp: "he" // PANIC: person struct does not have realProp property
}

```

## Structure methods
```
struct Person{
  name,
  age
}

Person -> sayHello(){
  var text = "Hi! im " + this.name
  print(text)
}

var person1 = Person{}
person1.name = "John"
person1.sayHello()
```

Evie does not have something like static method!

## Dictionaries
Dictionaries are similar to python dictionaries and javascript objects.
Trying to access or modify some prop that is not in the dictionary, will throw an error
```
var dict = {
  a: 1,
  b: 2
}

dict["a"] = 3
dict["c"] = 5 // PANIC: dict does not have key 'c'

dict.add("c", 4)
dict.remove("c")
dict.has("a")

print(dict["a"])

```
## Arrays
Arrays in Evie are similar to others languages like Python or JavaScript
```
var arr = [1, 2, "string", [3,4,5]]

arr[1] = 3
arr[10] = 5 // PANIC: Out Of Bounds!

arr.add(5)       // [1, 2, "string", [3,4,5], 5]
arr.addFirst(-2) // [-2, 1, 2, "string", [3,4,5], 5]
arr.has(2)       // true
arr.len()        // 6
arr.find("string") // 3

print(arr[2])

var slice = arr[1:-1]
var firstTwoValues = arr[:2]
```

## If statements
```
if 5 > 6{
  print("this is false")
}else if 5 > 8{
  print("this is false")
}else{
  print("any is true")
}
```
## Loops
```
var list = [1,2,3]

for item in list {
  print(item)
}

for index,item in list {
  print(index)
  print(item)
}

var dict = {a:1, b:2}

for key, val in dict{
  print("key: " + string(key))
  print("val: " + string(val))
}

var i = 0

loop {
  i = i + 1
  
  if i == 100{
    break
  }
  
  print("this will loop forever until i equals to 100")
}
```

Unlike many others languages, you CAN NOT redeclare variables inside loops
```
loop {
  var i = 0   // This will panic after the first iteration because 'i' will be redeclared
}

var i = 0

loop {
  i = 0   // This will loop forever, but will not give any errors
}

```
Also empty loops generates a panic
```
for item in items{

}   // PANIC: empty loop

loop {} // PANIC: empty loop

```

Documentation still in development