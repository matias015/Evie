# A simple scripting language

Evie is a personal project of a simple interpreted language with features like structures, dictionaries, modules, and more.

## Variables
You can not redeclare variables or access to non declarated ones
```
var x = 20
x = "Now x is a text"
```

## Functions
```
fn add(n1, n2){
  return n1 + n2
}

var result = add(2,4)
print(result)
```

## Structures
Like variables, you can not set, modify or access to a non defined property
```
struct Person{
  name
  age
}

var person1 = Person{}
person1.name = "John"
person1.age = 25

var person2 = Person{
  name: "Pedro"
  age: 24
}
```

## Structure methods
```
struct Person{
  name
  age
}

Person -> sayHello(){
  var text = "Hi, my name is: "
  print(text + this.name)
}

var person1 = Person{}
person1.name = "John"
person1.sayHello()
```

## Dictionaries
Dictionaries are similar to python dictionaries and javascript objects.
Trying to access or modify some prop that is not in the dictionary, will throw an error
```
var dict = {
  a: 1,
  b: 2
}

dict["a"] = 3

print(dict["a"])

```
## Arrays
Arrays in Evie are similar to others languages like Python or JavaScript
```
var arr = [1, 2, "string", [3,4,5]]

arr[1] = 3
arr[10] = 5 // Out Of Bounds!!

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

// PANIC: empty loops

for item in items{

}   

loop {}

```

## Built In Methods
```
input() // Captures and returns the user console input
print(...args) // Complex elements may not be printed correctly right now
number(arg)    // Parse the given value to a number
int(arg)       // Parse the given value to a number but also an integer
string(arg)    // Parse the given value to a string
bool(arg)      // Parse the given value to a bool
isNothing(...args) // Check if some of the given values are Nothing (equivalent to Null, None, Nill, etc)
type(arg)      // Return the type of the given value
time()         // The number of milliseconds since January 1, 1970
panic(msg)     // Throws an error with the given message
getArgs()      // Returns the execution arguments
```

Features and Documentation still in development